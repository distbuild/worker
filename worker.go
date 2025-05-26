package main

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"math"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"distbuild/boong/utils"
	"distbuild/boong/worker/proto"
)

var (
	BuildTime string
	CommitID  string
)

var (
	listenAddress string
	workSpacePath string
)

var rootCmd = &cobra.Command{
	Use:     "worker",
	Short:   "boong worker",
	Version: BuildTime + "-" + CommitID,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		if err := validArgs(ctx); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		if err := run(ctx); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	},
}

type server struct {
	proto.UnimplementedBuildServiceServer
}

// nolint:gochecknoinits
func init() {
	cobra.OnInitialize()

	rootCmd.PersistentFlags().StringVarP(&listenAddress, "listen-address", "l", ":39090", "listen address")
	rootCmd.PersistentFlags().StringVarP(&workSpacePath, "workspace-path", "w", "", "workspace path")

	_ = rootCmd.MarkFlagRequired("workspace-path")

	rootCmd.Root().CompletionOptions.DisableDefaultCmd = true
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func validArgs(_ context.Context) error {
	if listenAddress == "" {
		return errors.New("invalid listen address\n")
	}

	if workSpacePath == "" {
		return errors.New("workspace path is required\n")
	}

	if _, err := os.Stat(workSpacePath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return errors.New("workspace path does not exist\n")
		}
	}

	return nil
}

func run(_ context.Context) error {
	options := []grpc.ServerOption{grpc.MaxRecvMsgSize(math.MaxInt32), grpc.MaxSendMsgSize(math.MaxInt32)}

	s := grpc.NewServer(options...)
	proto.RegisterBuildServiceServer(s, &server{})

	l, _ := net.Listen("tcp", listenAddress)
	_ = s.Serve(l)

	return nil
}

func (s *server) SendBuild(stream proto.BuildService_SendBuildServer) error {
	var files []*proto.BuildFile
	var rule string
	var targets []string
	var id string
	var status bool

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		files = req.GetBuildFiles()
		rule = req.GetBuildRule()
		targets = req.GetBuildTargets()
		id = req.GetBuildID()
	}

	if err := s.validateBuildFile(files); err != nil {
		return errors.Wrap(err, "failed to validate build file\n")
	}

	if err := s.setBuildPath(targets); err != nil {
		return errors.Wrap(err, "failed to set build path\n")
	}

	t, err := s.runBuild(rule, targets, id)
	if err != nil {
		return errors.Wrap(err, "failed to run build\n")
	}

	if len(t) != 0 {
		status = true
	}

	if err := stream.Send(&proto.BuildReply{
		BuildTargets: t,
		BuildStatus:  status,
		BuildID:      id,
	}); err != nil {
		return errors.Wrap(err, "failed to send reply\n")
	}

	return nil
}

func (s *server) validateBuildFile(buildFiles []*proto.BuildFile) error {
	var _path string

	for _, file := range buildFiles {
		_path = filepath.Join(workSpacePath, file.FilePath)
		dir := filepath.Dir(_path)
		if _, err := os.Stat(_path); os.IsNotExist(err) {
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				if err := os.MkdirAll(dir, 0755); err != nil {
					return errors.Wrap(err, "failed to create directory\n")
				}
			}
		}

		if err := os.WriteFile(_path, file.GetFileData(), 0644); err != nil {
			return errors.Wrap(err, "failed to write file\n")
		}

		fileChecksum, err := utils.Checksum(_path)
		if err != nil {
			return errors.Wrap(err, "failed to calculate checksum\n")
		}
		if fileChecksum != file.GetCheckSum() {
			return errors.New("checksum mismatch\n")
		}
	}

	return nil
}

func (s *server) setBuildPath(targets []string) error {
	for _, target := range targets {
		if err := os.MkdirAll(filepath.Join(workSpacePath, filepath.Dir(target)), 0755); err != nil {
			return errors.Wrap(err, "failed to make directory\n")
		}
	}

	return nil
}

func (s *server) runBuild(rule string, targets []string, id string) ([]*proto.BuildTarget, error) {
	var buf []*proto.BuildTarget

	goos := runtime.GOOS

	args := strings.Fields(rule)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = workSpacePath

	if _, err := cmd.CombinedOutput(); err != nil {
		return nil, errors.Wrap(err, "failed to run command\n")
	}

	for _, target := range targets {
		p := filepath.Join(workSpacePath, target)
		if _, err := os.Stat(p); err != nil {
			if os.IsNotExist(err) {
				files, err := os.ReadDir(filepath.Dir(p))
				if err != nil {
					return nil, errors.Wrap(err, "failed to read directory\n")
				}
				for _, file := range files {
					if goos == "windows" {
						if filepath.Ext(file.Name()) == ".exe" {
							p = filepath.Join(filepath.Dir(p), filepath.Base(p)+".exe")
							break
						}
					}
				}
			} else {
				return nil, errors.Wrap(err, "failed to stat file\n")
			}
		}
		sum, err := utils.Checksum(p)
		if err != nil {
			return nil, errors.Wrap(err, "failed to calculate checksum\n")
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read file\n")
		}
		t := &proto.BuildTarget{
			TargetPath: target,
			Checksum:   sum,
			TargetData: data,
		}
		buf = append(buf, t)
	}

	return buf, nil
}
