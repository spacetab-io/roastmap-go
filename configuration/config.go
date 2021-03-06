package cfg

import (
	"fmt"
	"log"
	"strings"
	"time"

	config "github.com/spacetab-io/configuration-go"
	"gopkg.in/yaml.v2"

	"github.com/spacetab-io/prerender-go/pkg/models"
)

//Config config object
var Config *configuration

type serviceConfig struct {
	Name      string `yaml:"name"`
	About     string `yaml:"about"`
	Version   string `yaml:"version"`
	Docs      string `yaml:"docs"`
	Contacts  string `yaml:"contacts"`
	Copyright string `yaml:"copyright"`
}

type S3Config struct {
	Region       string `yaml:"region"`
	Bucket       string `yaml:"bucket"`
	BucketFolder string `yaml:"bucket_folder"`
	AccessKeyID  string `yaml:"access_key_id"`
	SecretKey    string `yaml:"secret_key"`
	CDNUrl       string `yaml:"cdn_url"`
}

type LocalStorageConfig struct {
	StoragePath string `yaml:"storage_path"`
}

type ElementConfig struct {
	Type      string `yaml:"type"`
	ID        string `yaml:"id"`
	Class     string `yaml:"class"`
	Attribute struct {
		Name  string `yaml:"name"`
		Value string `yaml:"value"`
	} `yaml:"attribute"`
}

type lookupConfig struct {
	Headless     bool     `yaml:"headless"`
	Type         string   `yaml:"type"`
	SitemapURLs  []string `yaml:"sitemaps"`
	PageURLs     []string `yaml:"urls"`
	BaseURL      string   `yaml:"base_url"`
	ParamsToSave []string `yaml:"get_params_to_save"`
}

func (c lookupConfig) GetSourceURL() string {
	switch c.Type {
	case models.LookupTypeSitemaps:
		return strings.Join(c.SitemapURLs, ",")
	case models.LookupTypeURLs:
		return strings.Join(c.PageURLs, ", ")
	}

	return ""
}

type viewportConfig struct {
	Width  int64 `yaml:"width"`
	Height int64 `yaml:"height"`
}

type PrerenderConfig struct {
	UserAgent       string         `yaml:"user_agent"`
	FilePostfix     string         `yaml:"file_postfix"`
	ConcurrentLimit int            `yaml:"concurrent_limit"`
	Lookup          lookupConfig   `yaml:"lookup"`
	WaitFor         string         `yaml:"wait_for"`
	ConsoleString   string         `yaml:"console_string"`
	SleepTime       time.Duration  `yaml:"sleep_time"`
	Element         ElementConfig  `yaml:"element"`
	Viewport        viewportConfig `yaml:"viewport"`
}

func (ec ElementConfig) GetWaitElement() string {
	elem := ec.Type

	if ec.ID != "" {
		elem = "#" + ec.ID
	}

	if ec.Class != "" && ec.ID == "" {
		elem = "." + ec.Class
	}

	return elem
}

func (ec ElementConfig) GetWaitElementAttr(attrValue string) string {
	elem := ec.GetWaitElement()

	if ec.Attribute.Name != "" && ec.Attribute.Value != "" {
		elem += fmt.Sprintf("[%s=%s]", ec.Attribute.Name, attrValue)
	}

	return elem
}

type StorageConfig struct {
	Type  string             `yaml:"type"`
	Local LocalStorageConfig `yaml:"local"`
	S3    S3Config           `yaml:"s3"`
}

type sitemapConfig struct {
	URL string `yaml:"url"`
}

// configuration file structure
type configuration struct {
	Service   serviceConfig   `yaml:"service"`
	Sitemap   sitemapConfig   `yaml:"sitemap"`
	Storage   StorageConfig   `yaml:"storage"`
	Prerender PrerenderConfig `yaml:"prerender"`
}

//Init initiates app config
func Init() error {
	configPath := config.GetEnv("CONFIG_PATH", "")
	configBytes, err := config.ReadConfigs(configPath)

	if err != nil {
		log.Printf("[config] read error: %+v", err)
		return err
	}

	if err = yaml.Unmarshal(configBytes, &Config); err != nil {
		log.Printf("[config] unmarshal error for bytes: %+v", configBytes)
		return err
	}

	return nil
}
