package configo

import (
	"github.com/hashicorp/consul/api"
	"strings"
)

type ConsulOptions struct {
	DefaultPath string
	Path string
	Prefix string
	Host string
	Token string
}

type Consul struct {
	options ConsulOptions
}

func NewConsulSource(co ConsulOptions) *Consul {
	c := &Consul{options: co}

	return c
}

func (c *Consul) GetVariables () (map[string]string, error) {
	return c.getConsulVariables()
}

func (c *Consul) getConsulVariables () (map[string]string, error) {
	vars := make(map[string]string)

	kv, err := c.getKVClient()
	if err != nil {
		return nil, err
	}

	for _, p := range []string{c.options.DefaultPath, c.options.Path} {
		if len(p) < 1 {
			continue
		}
		err = c.readVariables(p, kv, vars)
		if err != nil {
			return nil, err
		}
	}

	return vars, err
}

func (c *Consul) getKVClient () (*api.KV, error) {
	cfg := &api.Config{
		Address:    c.options.Host,
		Token:		c.options.Token,
	}
	client, err := api.NewClient(cfg)
	if err == nil {
		return client.KV(), err
	}

	return nil, err
}

func (c *Consul) readVariables (path string, kv *api.KV, m map[string]string) error {
	qo := &api.QueryOptions{}
	pairs, _, err := kv.List(path, qo)
	if err != nil {
		return err
	}
	for _, pair := range pairs{
		key := c.prepareKeyString(path, pair.Key)
		if len(key) > 0 {
			m[c.options.Prefix + key] = string(pair.Value)
		}
	}

	return nil
}

func (c *Consul) prepareKeyString(path, key string) string {
	k := strings.Replace(key, path, "", 1)
	k = strings.Trim(k, "/")
	k = strings.ReplaceAll(k, "/", "-")

	return k
}