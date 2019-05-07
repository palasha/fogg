package v2

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"reflect"

	"github.com/chanzuckerberg/fogg/config/v1"
	"github.com/chanzuckerberg/fogg/errs"
)

func ReadConfig(f io.Reader) (*Config, error) {
	c := &Config{
		Docker: false,
	}
	b, e := ioutil.ReadAll(f)
	if e != nil {
		return nil, errs.WrapUser(e, "unable to read config")
	}
	e = json.Unmarshal(b, c)

	return c, errs.WrapUser(e, "unable to parse config")
}

type Config struct {
	Accounts map[string]Account   `json:"accounts,omitempty"`
	Defaults Defaults             `json:"defaults" validate:"required"`
	Docker   bool                 `json:"docker,omitempty"`
	Envs     map[string]Env       `json:"envs,omitempty"`
	Global   Component            `json:"global,omitempty"`
	Modules  map[string]v1.Module `json:"modules,omitempty"`
	Plugins  v1.Plugins           `json:"plugins,omitempty"`
	Tools    Tools                `json:"tools,omitempty"`
	Version  int                  `json:"version" validate:"required,eq=2"`
}

type Common struct {
	Backend          Backend           `json:"backend,omitempty"`
	ExtraVars        map[string]string `json:"extra_vars,omitempty"`
	Owner            string            `json:"owner,omitempty" `
	Project          string            `json:"project,omitempty" `
	Providers        Providers         `json:"providers,omitempty" `
	TerraformVersion string            `json:"terraform_version,omitempty"`
}

type Defaults struct {
	Common
}

type Account struct {
	Common
}

type Tools struct {
	TravisCI *v1.TravisCI `json:"travis_ci,omitempty"`
	TfLint   *v1.TfLint   `json:"tflint,omitempty"`
}

type Env struct {
	Common

	Components map[string]Component `json:"components"`
}

type Component struct {
	Common

	EKS          *v1.EKSConfig     `json:"eks,omitempty"`
	Kind         *v1.ComponentKind `json:"kind,omitempty"`
	ModuleSource *string           `json:"module_source"`
}

type Providers struct {
	AWS *AWSProvider `json:"aws"`
}

type AWSProvider struct {
	// the aws provider is optional (above) but if supplied you must set account id and region
	AccountID         *int64   `json:"account_id" validate:"required"`
	AdditionalRegions []string `json:"additional_regions"`
	Profile           *string  `json:"profile"`
	Region            *string  `json:"region" validate:"required"`
	Version           *string  `json:"version,omitempty"`
}

type Backend struct {
	Bucket      string `json:"bucket,omitempty"`
	DynamoTable string `json:"dynamodb_table,omitempty"`
	Profile     string `json:"profile,omitempty"`
	Region      string `json:"region,omitempty"`
}

// Generate is used for test/quick integration. There are supposedly ways to do this without polluting the public
// api, but givent that fogg isn't an api, it doesn't seem like a big deal
func (c *Config) Generate(r *rand.Rand, size int) reflect.Value {
	// TODO write this to be part of tests https://github.com/shiwano/submarine/blob/5c02c0cfcf05126454568ef9624550eb0d84f86c/server/battle/src/battle/util/util_test.go#L19

	fmt.Println("generate")
	conf := &Config{}

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	randString := func(r *rand.Rand, n int) string {
		b := make([]byte, n)
		for i := range b {
			b[i] = letterBytes[r.Intn(len(letterBytes))]
		}
		return string(b)
	}

	randNonEmptyString := func(r *rand.Rand, s int) string {
		return "asdf"
	}

	randStringPtr := func(r *rand.Rand, s int) *string {
		str := randString(r, s)
		return &str
	}

	randStringMap := func(r *rand.Rand, s int) map[string]string {
		m := map[string]string{}

		for i := 0; i < s; i++ {
			m[randNonEmptyString(r, s)] = randString(r, s)
		}

		return map[string]string{}
	}

	randInt64Ptr := func(r *rand.Rand, s int) *int64 {
		if r.Float32() < 0.5 {
			i := r.Int63n(int64(size))
			return &i
		} else {
			var i *int64
			return i
		}
	}

	randAWSProvider := func(r *rand.Rand, s int) *AWSProvider {
		return &AWSProvider{
			AccountID: randInt64Ptr(r, size),
			Region:    randStringPtr(r, s),
			Profile:   randStringPtr(r, s),
			Version:   randStringPtr(r, s),
		}
	}

	// we treat these as opaque strings for now
	randVersion := func(r *rand.Rand, s int) string {
		return randString(r, s)
	}

	randCommon := func(r *rand.Rand, s int) Common {
		c := Common{
			Backend: Backend{
				Bucket: randString(r, s),
				Region: randString(r, s),
			},
			ExtraVars:        randStringMap(r, s),
			Owner:            randString(r, s),
			Project:          randString(r, s),
			Providers:        Providers{AWS: randAWSProvider(r, s)},
			TerraformVersion: randVersion(r, s),
		}
		return c
	}

	conf.Version = 2

	conf.Defaults = Defaults{
		Common: randCommon(r, size),
	}

	// tools

	conf.Tools = Tools{}

	if r.Float32() < 0.5 {
		conf.Tools.TravisCI = &v1.TravisCI{
			Enabled:     r.Float32() < 0.5,
			TestBuckets: r.Intn(size),
		}
	}

	conf.Accounts = map[string]Account{}
	acctN := r.Intn(size)

	for i := 0; i < acctN; i++ {
		acctName := randString(r, size)
		conf.Accounts[acctName] = Account{
			Common: randCommon(r, size),
		}

	}

	conf.Envs = map[string]Env{}
	envN := r.Intn(size)

	for i := 0; i < envN; i++ {
		envName := randString(r, size)
		e := Env{
			Common: randCommon(r, size),
		}
		e.Components = map[string]Component{}
		compN := r.Intn(size)

		for i := 0; i < compN; i++ {
			compName := randString(r, size)
			e.Components[compName] = Component{
				Common: randCommon(r, size),
			}
		}
		conf.Envs[envName] = e

	}

	return reflect.ValueOf(conf)
}