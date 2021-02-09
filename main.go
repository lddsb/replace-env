package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "replace-env"
	app.Usage = "A tiny tool for replace .env and .env.json in CI/CD"
	// Authors
	app.Authors = []*cli.Author{
		{
			Name:  "Roddy Happy",
			Email: "luodi0128@gmail.com",
		},
	}

	// Copyright
	year := time.Now().Year()
	var cpRight string
	if year == 2021 {
		cpRight = "2021"
	} else {
		cpRight = fmt.Sprintf("2021-%d", year)
	}
	app.Copyright = fmt.Sprintf("© %s Roddy Happy", cpRight)

	// Version
	app.Version = "0.1.210209"
	app.Compiled = time.Now()

	// Action
	app.Commands = command()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}

var branchEnv string

func command() []*cli.Command {
	return []*cli.Command{
		{
			Name:    "json",
			Aliases: []string{"j"},
			Action: func(c *cli.Context) error {
				sourceFile := c.Args().Get(0)
				if sourceFile == "" {
					log.Fatalln("source file cannot be null")
				}
				outputFile := c.Args().Get(1)
				jsonFile(sourceFile, outputFile)
				return nil
			},
			ArgsUsage: "replace-env j [--branch-env CI_COMMIT_BRANCH] source-file [output-file]",
			Flags:     commonFlags(),
		},
		{
			Name:    "env",
			Aliases: []string{"e"},
			Action: func(c *cli.Context) error {
				sourceFile := c.Args().Get(0)
				if sourceFile == "" {
					log.Fatalln("source file cannot be null")
				}
				outputFile := c.Args().Get(1)
				dotEnv(sourceFile, outputFile)
				return nil
			},
			ArgsUsage: "replace-env e [--branch-env CI_COMMIT_BRANCH] source-file [output-file]",
			Flags:     commonFlags(),
		},
	}
}

// commonFlags
func commonFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "branch-env",
			Usage:       "change branch environment",
			Value:       "CI_COMMIT_BRANCH",
			Destination: &branchEnv,
		},
	}
}

// jsonFile replace json file, then output to a new file
func jsonFile(source string, outfile string) {
	jsonByte, err := ioutil.ReadFile(source)
	if err != nil {
		log.Fatalln(err)
		return
	}
	var mapResult map[string]interface{}
	if err := json.Unmarshal(jsonByte, &mapResult); err != nil {
		log.Fatalln(err)
		return
	}

	mp := jsonRecursive(mapResult)
	// pretty JSON marshal
	jsonByte, err = json.MarshalIndent(mp, "", "	")
	if outfile == "" {
		fmt.Println(string(jsonByte))
	} else {
		err = ioutil.WriteFile(outfile, jsonByte, os.FileMode(0644))
		if err != nil {
			log.Fatalln(err)
			return
		}
	}
}

// jsonRecursive recursively traverse the JSON file
func jsonRecursive(mp map[string]interface{}) map[string]interface{} {
	for k, v := range mp {
		env := os.Getenv(branchEnv + strings.ToUpper(k))
		switch x := v.(type) {
		case map[string]interface{}:
			mp[k] = jsonRecursive(x)
		case float64:
			if env == "" {
				env = "0"
			}
			envFloat, err := strconv.ParseFloat(env, 64)
			if err != nil {
				log.Fatalln(err)
			}

			mp[k] = envFloat
		case bool:
			if strings.ToLower(env) != "false" && strings.ToLower(env) != "true" {
				env = "false"
			}

			envBool, err := strconv.ParseBool(env)
			if err != nil {
				log.Fatalln(err)
			}
			mp[k] = envBool
		default:
			mp[k] = env
		}
	}

	return mp
}

// dotEnv replace .env file
func dotEnv(sourceFile string, outputFile string) {
	envBytes, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		log.Fatalln(err)
	}
	envSlice := strings.Split(string(envBytes), "\n")
	newStringSlice := []string{}
	for _, v := range envSlice {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}

		env := strings.SplitN(v, "=", 2)
		newStringSlice = append(newStringSlice, fmt.Sprintf("%s=%v", env[0], os.Getenv(branchEnv+strings.ToUpper(env[0]))))
	}

	if outputFile == "" {
		for _, lineStr := range newStringSlice {
			fmt.Println(lineStr)
		}
	} else {
		ioutil.WriteFile(outputFile, []byte(strings.Join(newStringSlice, "\n")), 0644)
		if err != nil {
			log.Fatalln(err)
			return
		}
	}
}
