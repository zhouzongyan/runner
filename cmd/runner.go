package cmd

import (
	"fmt"
	"github.com/gokins-main/core/utils"
	"github.com/gokins-main/runner/runners"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"github.com/sirupsen/logrus"
	"io"
	"strconv"
)

type HbtpRunner struct {
	cfg Config
}

func (c *HbtpRunner) ServerInfo() runners.ServerInfo {
	info := &runners.ServerInfo{}
	err := c.doHbtpJson("ServerInfo", nil, info)
	if err != nil {
		logrus.Errorf("hbtp ServerInfo err:%v", err)
	}
	return *info
}
func (c *HbtpRunner) PullJob(plugs []string) (*runners.RunJob, error) {
	rt := &runners.RunJob{}
	err := c.doHbtpJson("PullJob", plugs, rt)
	return rt, err
}
func (c *HbtpRunner) CheckCancel(buildId string) bool {
	code, bts, err := c.doHbtpString("CheckCancel", nil, hbtp.Map{
		"buildId": buildId,
	})
	rts := string(bts)
	if err != nil || code != hbtp.ResStatusOk {
		return false
	}
	return rts == "true"
}
func (c *HbtpRunner) Update(m *runners.UpdateJobInfo) error {
	code, bts, err := c.doHbtpString("Update", m)
	if err != nil {
		return err
	}
	if code != hbtp.ResStatusOk {
		return fmt.Errorf("%s", string(bts))
	}
	return nil
}
func (c *HbtpRunner) UpdateCmd(buildId, jobId, cmdId string, fs, codes int) error {
	code, bts, err := c.doHbtpString("UpdateCmd", nil, hbtp.Map{
		"buildId": buildId,
		"jobId":   jobId,
		"cmdId":   cmdId,
		"fs":      fs,
		"code":    codes,
	})
	if err != nil {
		return err
	}
	if code != hbtp.ResStatusOk {
		return fmt.Errorf("%s", string(bts))
	}
	return nil
}
func (c *HbtpRunner) PushOutLine(buildId, jobId, cmdId, bs string, iserr bool) error {
	code, bts, err := c.doHbtpString("PushOutLine", nil, hbtp.Map{
		"buildId": buildId,
		"jobId":   jobId,
		"cmdId":   cmdId,
		"bs":      bs,
		"iserr":   iserr,
	})
	if err != nil {
		return err
	}
	if code != hbtp.ResStatusOk {
		return fmt.Errorf("%s", string(bts))
	}
	return nil
}
func (c *HbtpRunner) FindJobId(buildId, stgNm, stpNm string) (string, bool) {
	code, bts, err := c.doHbtpString("FindJobId", nil, hbtp.Map{
		"buildId": buildId,
	})
	rts := string(bts)
	if err != nil {
		return "", false
	}
	if code != hbtp.ResStatusOk {
		return "", false
	}
	return rts, true
}
func (c *HbtpRunner) ReadDir(fs int, buildId string, pth string) ([]*runners.DirEntry, error) {
	var rts []*runners.DirEntry
	err := c.doHbtpJson("ReadDir", nil, &rts, hbtp.Map{
		"buildId": buildId,
		"pth":     pth,
		"fs":      fs,
	})
	if err != nil {
		return nil, err
	}
	return rts, nil
}
func (c *HbtpRunner) ReadFile(fs int, buildId string, pth string) (int64, io.ReadCloser, error) {
	req := c.newHbtpReq("ReadFile")
	req.ReqHeader().Set("buildId", buildId)
	req.ReqHeader().Set("pth", pth)
	req.ReqHeader().Set("fs", fs)
	err := req.Do(nil, nil)
	if err != nil {
		return 0, nil, err
	}
	defer req.Close()
	rs := string(req.ResBodyBytes())
	if req.ResCode() != hbtp.ResStatusOk {
		return 0, nil, fmt.Errorf("%s", rs)
	}
	sz, err := strconv.ParseInt(rs, 10, 64)
	if err != nil {
		return 0, nil, err
	}
	return sz, req.Conn(true), nil
}
func (c *HbtpRunner) GetEnv(buildId, jobId, key string) (string, bool) {
	code, bts, err := c.doHbtpString("GetEnv", nil, hbtp.Map{
		"buildId": buildId,
		"jobId":   jobId,
		"key":     key,
	})
	rts := string(bts)
	if err != nil {
		return "", false
	}
	if code != hbtp.ResStatusOk {
		return "", false
	}
	return rts, true
}
func (c *HbtpRunner) GenEnv(buildId, jobId string, env utils.EnvVal) error {
	code, bts, err := c.doHbtpString("GenEnv", env, hbtp.Map{
		"buildId": buildId,
		"jobId":   jobId,
	})
	if err != nil {
		return err
	}
	if code != hbtp.ResStatusOk {
		return fmt.Errorf("%s", string(bts))
	}
	return nil
}
func (c *HbtpRunner) UploadFile(fs int, buildId, jobId string, dir, pth string) (io.WriteCloser, error) {
	req := c.newHbtpReq("UploadFile")
	req.ReqHeader().Set("buildId", buildId)
	req.ReqHeader().Set("jobId", jobId)
	req.ReqHeader().Set("dir", dir)
	req.ReqHeader().Set("pth", pth)
	req.ReqHeader().Set("fs", fs)
	err := req.Do(nil, nil)
	if err != nil {
		return nil, err
	}
	defer req.Close()
	rs := string(req.ResBodyBytes())
	if req.ResCode() != hbtp.ResStatusOk {
		return nil, fmt.Errorf("%s", rs)
	}
	return req.Conn(true), nil
}
func (c *HbtpRunner) FindArtVersionId(buildId, idnt string, name string) (string, error) {
	code, bts, err := c.doHbtpString("FindArtVersionId", nil, hbtp.Map{
		"buildId": buildId,
		"idnt":    idnt,
		"name":    name,
	})
	rts := string(bts)
	if err != nil {
		return "", err
	}
	if code != hbtp.ResStatusOk {
		return "", fmt.Errorf("%s", string(bts))
	}
	return rts, nil
}
func (c *HbtpRunner) NewArtVersionId(buildId, idnt string, name string) (string, error) {
	code, bts, err := c.doHbtpString("NewArtVersionId", nil, hbtp.Map{
		"buildId": buildId,
		"idnt":    idnt,
		"name":    name,
	})
	rts := string(bts)
	if err != nil {
		return "", err
	}
	if code != hbtp.ResStatusOk {
		return "", fmt.Errorf("%s", string(bts))
	}
	return rts, nil
}
