package infra

import (
	"archive/tar"
	"bytes"
	"context"
	"log"
	"strings"
	"time"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/zhohov/contest-service/judge-service/internal/domain"
)

type ExecCodeResult struct {
	stdout string
	stderr string
	time   float64
}

type DockerJudger struct {
	client *client.Client
}

func NewDockerJudger(client *client.Client) *DockerJudger {
	return &DockerJudger{
		client: client,
	}
}

func (e *DockerJudger) createContainer(ctx context.Context) (client.ContainerCreateResult, error) {
	container, err := e.client.ContainerCreate(ctx, client.ContainerCreateOptions{
		Image: "python:3.11-slim",
		Config: &container.Config{
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			Tty:          false,
			OpenStdin:    true,
		},
	})
	if err != nil {
		return client.ContainerCreateResult{}, err
	}

	if _, err := e.client.ContainerStart(ctx, container.ID, client.ContainerStartOptions{}); err != nil {
		return client.ContainerCreateResult{}, err
	}

	return container, nil
}

func (e *DockerJudger) writeCode(ctx context.Context, containerID string, code string) error {
	content := []byte(code)

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	hdr := &tar.Header{
		Name: "main.py",
		Mode: 0644,
		Size: int64(len(content)),
	}
	tw.WriteHeader(hdr)
	tw.Write(content)
	tw.Close()

	_, err := e.client.CopyToContainer(ctx, containerID, client.CopyToContainerOptions{
		DestinationPath: "/tmp",
		Content:         buf,
	})
	if err != nil {
		log.Println(err)
	}

	return nil
}

func (e *DockerJudger) executeCode(ctx context.Context, containerID string, stdin string) (ExecCodeResult, error) {
	start := time.Now()

	exec, err := e.client.ExecCreate(ctx, containerID, client.ExecCreateOptions{
		Cmd:          []string{"python", "/tmp/main.py"},
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return ExecCodeResult{}, err
	}

	attach, err := e.client.ExecAttach(ctx, exec.ID, client.ExecAttachOptions{})
	if err != nil {
		return ExecCodeResult{}, err
	}
	defer attach.Close()

	_, _ = attach.Conn.Write([]byte(stdin + "\n"))
	attach.CloseWrite()

	var stdoutBuf, stderrBuf bytes.Buffer
	_, err = stdcopy.StdCopy(&stdoutBuf, &stderrBuf, attach.Reader)
	if err != nil {
		return ExecCodeResult{}, err
	}

	return ExecCodeResult{
		stdout: stdoutBuf.String(),
		stderr: stderrBuf.String(),
		time:   time.Since(start).Seconds(),
	}, nil
}

func (e *DockerJudger) Execute(ctx context.Context, code string, inputs []string) ([]ExecCodeResult, error) {
	container, err := e.createContainer(ctx)
	if err != nil {
		return []ExecCodeResult{}, err
	}
	defer e.client.ContainerRemove(ctx, container.ID, client.ContainerRemoveOptions{Force: true})

	if err := e.writeCode(ctx, container.ID, code); err != nil {
		return []ExecCodeResult{}, err
	}

	res := make([]ExecCodeResult, 0, len(inputs))
	for _, input := range inputs {
		execRes, err := e.executeCode(ctx, container.ID, input)
		if err != nil {
			log.Println(err)
		}
		res = append(res, execRes)
	}

	return res, nil
}

type SubmitTask struct {
	Stdout string  `json:"stdout"`
	Stderr string  `json:"stderr"`
	Time   float64 `json:"time"`
	Pass   bool    `json:"pass"`
}

type SubmitResult struct {
	Tests []SubmitTask `json:"tests"`
	Pass  bool         `json:"pass"`
}

type TestCase struct {
	Input    string
	Expected string
}

func (e *DockerJudger) Submit(ctx context.Context, code string, testCases []TestCase) (*domain.Submission, error) {
	inputs := make([]string, 0, len(testCases))
	for _, tc := range testCases {
		inputs = append(inputs, tc.Input)
	}
	executes, err := e.Execute(ctx, code, inputs)
	if err != nil {
		return &domain.Submission{}, err
	}

	result := domain.SubmissionResultAccepted
	for i, exec := range executes {
		if strings.TrimSpace(exec.stderr) != "" {
			result = domain.SubmissionResultRuntimeError
			break
		}
		if strings.TrimSpace(exec.stdout) != strings.TrimSpace(testCases[i].Expected) {
			result = domain.SubmissionResultWrongAnswer
			break
		}
	}

	return &domain.Submission{
		ID:     0,
		Result: result,
	}, nil
}
