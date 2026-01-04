package ossc

import "github.com/aliyun/aliyun-oss-go-sdk/oss"

var (
	client       *Client
	publicClient *Client
)

type Cfg struct {
	PublicEndpoint  string `mapstructure:"public_endpoint"`
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyId     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	UserFileBucket  string `mapstructure:"user_file_bucket"`
}

type Client struct {
	client *oss.Client
	bucket *oss.Bucket
}

func Init(cfg Cfg) (err error) {
	cli, err := oss.New(cfg.Endpoint, cfg.AccessKeyId, cfg.AccessKeySecret)
	if err != nil {
		return
	}
	userFileBucket, err := cli.Bucket(cfg.UserFileBucket)
	if err != nil {
		return
	}
	client = &Client{
		client: cli,
		bucket: userFileBucket,
	}

	publicCli, err := oss.New(cfg.PublicEndpoint, cfg.AccessKeyId, cfg.AccessKeySecret)
	if err != nil {
		return
	}
	publicBucket, err := publicCli.Bucket(cfg.UserFileBucket)
	if err != nil {
		return
	}
	publicClient = &Client{
		client: publicCli,
		bucket: publicBucket,
	}
	return
}

func Get() *Client {
	return client
}

func GetPublic() *Client {
	return publicClient
}

func (c *Client) UserFileBucket() *oss.Bucket {
	return c.bucket
}
