package configo

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func consulHandlers() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/v1/kv/ump/site20/application", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"LockIndex": 0,"Key": "ump/site20/application/","Flags": 0,"Value": null,"CreateIndex": 4652,"ModifyIndex": 4652},{"LockIndex": 0,"Key": "ump/site20/application/param1-text","Flags": 0,"Value": "c29tZXRleHQ=","CreateIndex": 4655,"ModifyIndex": 4655},{"LockIndex": 0,"Key": "ump/site20/application/param2-int","Flags": 0,"Value": "OQ==","CreateIndex": 4658,"ModifyIndex": 4658},{"LockIndex": 0,"Key": "ump/site20/application/param3-empty","Flags": 0,"Value": null,"CreateIndex": 4659,"ModifyIndex": 4659},{"LockIndex": 0,"Key": "ump/site20/application/param4-text-override","Flags": 0,"Value": "b3ZlcnJpZGU=","CreateIndex": 4661,"ModifyIndex": 4661},{"LockIndex": 0,"Key": "ump/site20/application/param5-folder/","Flags": 0,"Value": null,"CreateIndex": 4667,"ModifyIndex": 4667},{"LockIndex": 0,"Key": "ump/site20/application/param5-folder/subfolder/val-as-string","Flags": 0,"Value": "NTU1","CreateIndex": 4671,"ModifyIndex": 4671},{"LockIndex": 0,"Key": "ump/site20/application/param5-int-overide","Flags": 0,"Value": "Mw==","CreateIndex": 4664,"ModifyIndex": 4664}]`))
	})

	r.HandleFunc("/v1/kv/ump/site20/site2-test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"LockIndex": 0,"Key": "ump/site20/site2-test/","Flags": 0,"Value": null,"CreateIndex": 4683,"ModifyIndex": 4683},{"LockIndex": 0,"Key": "ump/site20/site2-test/param4-text-override","Flags": 0,"Value": "c2l0ZTItdGVzdC50ZXh0","CreateIndex": 4687,"ModifyIndex": 4697},{"LockIndex": 0,"Key": "ump/site20/site2-test/param5-int-override","Flags": 0,"Value": "NDQ0","CreateIndex": 4689,"ModifyIndex": 4689}]`))
	})

	return r
}

func TestConsul(t *testing.T) {
	srv := httptest.NewServer(consulHandlers())
	defer srv.Close()

	co := ConsulOptions{
		DefaultPath: "ump/site20/application",
		Path:        "ump/site20/site2-test",
		Host:        srv.URL,
	}
	cs := NewConsulSource(co)
	LoadConfigs(cs)

	os.Setenv("test-env", "text")
	testEnv := EnvString("test-env", "")
	if testEnv != "text" {
		t.Errorf("Get from simple env error - need text, got %s", testEnv)
	}
	_, ok := ConfigVariables["test-env"]
	if !ok {
		t.Errorf("Caching simple env error")
	}

	asString := EnvString("param5-folder-subfolder-val-as-string", "empty")
	if asString != "555" {
		t.Errorf("Get integer as string error need 555, got %s", asString)
	}

	asInt := EnvInt("param5-folder-subfolder-val-as-string", 0)
	if asInt != 555 {
		t.Errorf("Get integer as string error need 555, got %d", asInt)
	}

	emptyString := EnvString("param3-empty", "not empty")
	if emptyString != "" {
		t.Errorf(`Get integer as string error need "", got %s`, emptyString)
	}

	emptyInt := EnvInt("param3-empty", 10)
	if emptyInt != 10 {
		t.Errorf(`Get integer as string error need 10, got %d`, emptyInt)
	}

	notExistParamString := EnvString("not-exist-param-string", "defaultString")
	if notExistParamString != "defaultString" {
		t.Errorf(`Get integer as string error need "defaultString", got %s`, notExistParamString)
	}

	notExistParamInt := EnvInt("not-exist-param-string", 30)
	if notExistParamInt != 30 {
		t.Errorf(`Get integer as string error need 30, got %d`, notExistParamInt)
	}

	overrideString := EnvString("param4-text-override", "defaultString")
	if overrideString != "site2-test.text" {
		t.Errorf(`Get integer as string error need "site2-test.text", got %s`, overrideString)
	}

	overrideInt := EnvInt("param5-int-override", 30)
	if overrideInt != 444 {
		t.Errorf(`Get integer as string error need 444, got %d`, overrideInt)
	}

	errorInt := EnvInt("param4-text-override", 30)
	if errorInt != 30 {
		t.Errorf(`Get integer as string error need 30, got %d`, errorInt)
	}
}
