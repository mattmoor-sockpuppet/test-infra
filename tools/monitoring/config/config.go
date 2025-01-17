/*
Copyright 2019 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// config is responsible for fetching, parsing config yaml file. It also allows user to
// retrieve a particular record from the yaml.

package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"gopkg.in/yaml.v2"
)

type alertCondition struct {
	JobNameRegex string `yaml:"job-name-regex"`
	Occurrences  int
	JobsAffected int `yaml:"jobs-affected"`
	PrsAffected  int `yaml:"prs-affected"`
	Period       int
}

type patternSpec struct {
	ErrorPattern string `yaml:"error-pattern"`
	Hint         string
	Alerts       []alertCondition
}

// Config stores all information read from the config yaml
type Config struct {
	Spec []patternSpec `yaml:"spec"`
}

// SelectedConfig stores the recovery hint as well as alert conditions for a selected error pattern
// and qualifying job name
type SelectedConfig struct {
	Hint         string
	Occurrences  int
	JobsAffected int
	PrsAffected  int
	Period       int
}

// applyDefaults set fields to desired defaults values if they are missing from yaml
func (selected SelectedConfig) applyDefaults() {
	if selected.Occurrences == 0 {
		selected.Occurrences = 1
	}
	if selected.JobsAffected == 0 {
		selected.JobsAffected = 1
	}
	if selected.PrsAffected == 0 {
		selected.PrsAffected = 1
	}
	if selected.Period == 0 {
		selected.Period = 24 * 60
	}
}

// Select gets the spec for a particular error pattern and a matching job name pattern
func (config Config) Select(pattern, jobName string) (SelectedConfig, error) {
	output := SelectedConfig{}
	noMatchError := fmt.Errorf("no spec found for pattern[%s] and jobName[%s]",
		pattern, jobName)
	for _, patternSpec := range config.Spec {
		if pattern == patternSpec.ErrorPattern {
			noMatchError = fmt.Errorf("spec found for pattern[%s], but no match for job name[%s]", pattern, jobName)
			output.Hint = patternSpec.Hint
			for _, alertCondition := range patternSpec.Alerts {
				matched, err := regexp.MatchString(alertCondition.JobNameRegex, jobName)
				if err != nil {
					log.Printf("Error matching pattern '%s' on string '%s': %v",
						alertCondition.JobNameRegex, jobName, err)
					continue
				}
				if matched {
					noMatchError = nil
					output.JobsAffected = alertCondition.JobsAffected
					output.Occurrences = alertCondition.Occurrences
					output.PrsAffected = alertCondition.PrsAffected
					output.Period = alertCondition.Period
					break
				}
			}
			break
		}
	}
	return output, noMatchError
}

// CollectErrorPatterns collects and returns all error patterns in the yaml file
func (config Config) CollectErrorPatterns() []string {
	var patterns []string
	for _, patternSpec := range config.Spec {
		patterns = append(patterns, patternSpec.ErrorPattern)
	}
	return patterns
}

// GetFileBytes retrieves a file by URL and returns its text content
func GetFileBytes(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

// CompilePatterns compiles the patterns from string to Regexp. In addition it returns the list of
// patterns that cannot be compiled
func CompilePatterns(patterns []string) ([]regexp.Regexp, []string) {
	var regexps []regexp.Regexp
	var badPatterns []string // patterns that cannot be compiled into regex

	for _, pattern := range patterns {
		r, err := regexp.Compile(pattern)
		if err != nil {
			log.Printf("Error compiling pattern [%s]: %v", pattern, err)
			badPatterns = append(badPatterns, pattern)
		} else {
			regexps = append(regexps, *r)
		}
	}
	return regexps, badPatterns
}

// ParseYaml reads the yaml text and converts it to the Config struct defined
func ParseYaml(url string) (*Config, error) {
	content, err := GetFileBytes(url)
	if err != nil {
		return nil, err
	}
	return newConfig(content)
}

func newConfig(text []byte) (*Config, error) {
	file := new(Config)
	if err := yaml.UnmarshalStrict(text, &file); err != nil {
		return file, err
	}
	return file, nil
}

// GetAllPatterns collects all regexp patterns, including both error message patterns
// and job name patterns
func (config *Config) GetAllPatterns() []string {
	var patterns []string
	for _, patternSpec := range config.Spec {
		patterns = append(patterns, patternSpec.ErrorPattern)
		for _, alertCondition := range patternSpec.Alerts {
			patterns = append(patterns, alertCondition.JobNameRegex)
		}
	}

	return patterns
}
