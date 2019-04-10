package ecsupdatenotify

import (
	"github.com/BurntSushi/toml"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"go.uber.org/zap"
)

const (
	ColorORANGE = "#F6D64F"
	ColorRED    = "#F08080"
)

var (
	pid                  = "/tmp/ecs-update-notify.pid"
	log                  Logger
	CheckFailureInterval = 20
)

// Config ... config.toml
type Config struct {
	Interval             int        `toml:"interval"`
	CheckFailureInterval int        `toml:"check_failure_interval"`
	Monitors             []*Monitor `toml:"monitor"`
}

// Monitor ... set from config.toml
type Monitor struct {
	Name            string `toml:"name"`
	AWSProfile      string `toml:"aws_profile"`
	AWSRegion       string `toml:"aws_region"`
	IncomingWebhook string `toml:"incoming_webhook"`
	Tasks           []*ECSTask
	Client          *Client
}

// Client ... Store ECS client with a session
type Client struct {
	ecs ecsiface.ECSAPI
}

// Logger ... Store logging
type Logger struct {
	sugar *zap.SugaredLogger
}

// ECSTask ...
type ECSTask struct {
	SlackMessage
	Name              string
	CurrRevision      int
	NextRevision      int
	isCurrReivision   bool
	TaskDefinitionArn *string
	FailureCount      int
}

// ECSFailureTask ...
type ECSFailureTask struct {
	SlackMessage
}

// Slack ... Store Attachment information
type Slack struct {
	Attachments []*Attachment `json:"attachments"`
}

// Attachment ... Slack Attachment Data
type Attachment struct {
	Color  string `json:"color,omitempty"`
	Title  string `json:"title,omitempty"`
	Text   string `json:"text,omitempty"`
	Footer string `json:"footer,omitempty"`
}

// LoadToml ... Load toml file
func LoadToml(c string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(c, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
