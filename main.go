package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type config struct {
	Language   string   `yaml:"language" json:"language"`
	Includes   []string `yaml:"includes" json:"includes"`
	Tags       []string `yaml:"tags" json:"tags"`
	Inputs     []string `yaml:"input" json:"input"`
	Output     string   `yaml:"output" json:"output"`
	Module     string   `yaml:"module" json:"module"`
	CopyReadme bool     `yaml:"copyReadme" json:"copyReadme"`

	commands []string
}

func main() {
	c := new(config)

	c.parse()
	c.validate()
	c.generateCommands()

	if err := c.run(); nil != err {
		panic(err)
	}
}

func (c *config) parse() {
	c.Language = os.Getenv("PLUGIN_LANGUAGE")

	includes := os.Getenv("PLUGIN_INCLUDES")
	if includes != "" {
		c.Includes = strings.Split(includes, ",")
		for i := range c.Includes {
			c.Includes[i] = "--proto_path=" + c.Includes[i]
		}
	}

	tags := os.Getenv("PLUGIN_TAGS")
	if tags != "" {
		c.Tags = strings.Split(tags, ",")
	}

	input := os.Getenv("PLUGIN_INPUTS")
	if input != "" {
		c.Inputs = strings.Split(input, ",")
	}

	c.Output = os.Getenv("PLUGIN_OUTPUT")
	c.Module = os.Getenv("PLUGIN_MODULE")
	c.CopyReadme, _ = strconv.ParseBool(os.Getenv("PLUGIN_COPYREADME"))
}

func (c *config) validate() {
	switch strings.ToUpper(c.Language) {
	case "GO", "DART", "JAVASCRIPT":
	default:
		panic("unsupported language.")
	}

	if len(c.Inputs) == 0 {
		panic("no input files.")
	}

	if c.Output == "" {
		panic("invalid output file.")
	}
}

func (c *config) generateCommands() {
	includes := strings.Join(c.Includes, ` `)
	tags := strings.Join(c.Tags, ` `)

	var commands []string

	switch strings.ToUpper(c.Language) {
	case "GO":
		pre := fmt.Sprintf(`protoc %s %s --go_out=plugins=grpc:%s`, includes, tags, c.Output)
		if c.Module != "" {
			pre += ` --go_opt=module=` + c.Module
		}

		for _, in := range c.Inputs {
			commands = append(commands, pre+` `+in)
		}
	case "DART":
		commands = append(commands, fmt.Sprintf(`protoc %s --dart_out=generate_kythe_info:%s/lib /usr/include/google/protobuf/*.proto`, tags, c.Output))

		for _, in := range c.Inputs {
			commands = append(
				commands,
				fmt.Sprintf(`cd lib/core && protoc %s --dart_out=generate_kythe_info:../../%s/lib %s && cd ../../`, tags, c.Output, in),
				fmt.Sprintf(`protoc %s %s --dart_out=generate_kythe_info:%s/lib %s`, includes, tags, c.Output, in),
			)
		}
	case "JAVASCRIPT":
		pre := fmt.Sprintf(`protoc %s %s --js_out=binary:%s`, includes, tags, c.Output)

		for _, in := range c.Inputs {
			commands = append(commands, pre+` `+in)
		}
	}

	if c.CopyReadme {
		commands = append(commands, fmt.Sprintf(`cp -f README.md %s/README.md`, c.Output))
	}

	c.commands = commands
}

func (c *config) run() (err error) {
	for _, command := range c.commands {
		cmd := exec.Command("sh", "-c", command)

		fmt.Println("--", command)

		var stdout io.ReadCloser
		if stdout, err = cmd.StdoutPipe(); err != nil {
			return
		}

		if err = cmd.Start(); nil != err {
			return
		}

		var line []byte
		if line, err = ioutil.ReadAll(stdout); nil == err {
			fmt.Println(string(line))
		}

		if err = cmd.Wait(); nil != err {
			return
		}

		stdout.Close()
	}

	return
}
