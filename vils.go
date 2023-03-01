package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
)

func main() {
	manifest, err := createManifest(os.Args[1:])
	fatal(err)

	manifestPath, err := writeManifest(manifest)
	fatal(err)
	defer os.Remove(manifestPath)

	fatal(editManifest(manifestPath))

	changedManifest, err := readManifest(manifestPath)
	fatal(err)

	ops := generateOperations(manifest, changedManifest)

	for _, op := range ops {
		fmt.Println(op.String())
		fatal(op.Apply())
	}
}

func generateOperations(manifest []string, changedManifest map[int]string) []operation {
	out := make([]operation, 0)

	for idx, src := range manifest {
		dst, exists := changedManifest[idx]
		if !exists {
			out = append(out, &remove{src})
		} else if src != dst {
			out = append(out, &rename{src, dst})
		}
	}
	return out
}
func readManifest(manifestPath string) (map[int]string, error) {
	pattern := regexp.MustCompile("^(\\d+)\t(.*)$")
	f, err := os.Open(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("reading manifest: %w", err)
	}
	defer f.Close()

	out := make(map[int]string)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		res := pattern.FindStringSubmatch(line)
		if res == nil {
			return nil, fmt.Errorf("reading manifest: garbage line '%s'", line)
		}

		idx, err := strconv.Atoi(res[1])
		if err != nil {
			return nil, fmt.Errorf("reading manifest: invalid integer '%s'", res[1])
		}
		out[idx] = res[2]
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading manifest: %w", err)
	}
	return out, nil
}

func editManifest(manifestPath string) error {
	cmd := exec.Command(getEditor(), manifestPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("editing manifest: %w", err)
	}
	return nil
}

func createManifest(arguments []string) ([]string, error) {
	paths := make([]string, 0, 10000)
	if len(arguments) == 0 {
		var err error
		paths, err = filepath.Glob("*")
		if err != nil {
			return nil, fmt.Errorf("listing arguments: %w", err)
		}
	} else {
		for _, path := range arguments {
			s, err := os.Stat(path)
			if err != nil {
				return nil, fmt.Errorf("listing arguments: %w", err)
			}

			if s.IsDir() {
				paths, err = appendDir(paths, path)
				if err != nil {
					return nil, fmt.Errorf("listing arguments: %w", err)
				}
			} else {
				paths = append(paths, path)
			}
		}
	}
	return paths, nil
}

func writeManifest(manifest []string) (string, error) {
	tmpf, err := os.CreateTemp("", "vils")
	if err != nil {
		return "", fmt.Errorf("writing manifest: %w", err)
	}
	for i, fn := range manifest {
		if _, err := fmt.Fprintf(tmpf, "%d\t%s\n", i, fn); err != nil {
			return "", fmt.Errorf("writing manifest: %w", err)
		}
	}
	tmpfName := tmpf.Name()
	tmpf.Close()
	return tmpfName, nil
}

func appendDir(paths []string, path string) ([]string, error) {
	infos, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, info := range infos {
		paths = append(paths, filepath.Join(path, info.Name()))
	}
	return paths, nil
}

func getEditor() string {
	if ed := os.Getenv("VISUAL"); ed != "" {
		return ed
	}
	if ed := os.Getenv("EDITOR"); ed != "" {
		return ed
	}
	return "vi"
}

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}
