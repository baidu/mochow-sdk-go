package api

import (
	"github.com/baidu/mochow-sdk-go/v2/client"
	"github.com/baidu/mochow-sdk-go/v2/http"
	"github.com/bytedance/sonic"
)

func CreateRole(cli client.Client, args *CreateRoleArgs) error {
	req := &client.BceRequest{}
	req.SetURI(getRoleURI())
	req.SetMethod(http.Post)
	req.SetParam("create", "")
	jsonBytes, err := sonic.Marshal(args)
	if err != nil {
		return err
	}
	body, err := client.NewBodyFromBytes(jsonBytes)
	if err != nil {
		return err
	}
	req.SetBody(body)

	resp := &client.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return err
	}
	defer resp.Body().Close()
	if resp.IsFail() {
		return resp.ServiceError()
	}
	return nil
}

func DropRole(cli client.Client, args *DropRoleArgs) error {
	req := &client.BceRequest{}
	req.SetURI(getRoleURI())
	req.SetMethod(http.Post)
	req.SetParam("drop", "")
	jsonBytes, err := sonic.Marshal(args)
	if err != nil {
		return err
	}
	body, err := client.NewBodyFromBytes(jsonBytes)
	if err != nil {
		return err
	}
	req.SetBody(body)

	resp := &client.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return err
	}
	defer resp.Body().Close()
	if resp.IsFail() {
		return resp.ServiceError()
	}
	return nil
}

func GrantRolePrivileges(cli client.Client, args *GrantRolePrivilegesArgs) error {
	req := &client.BceRequest{}
	req.SetURI(getRoleURI())
	req.SetMethod(http.Post)
	req.SetParam("grantPrivileges", "")
	jsonBytes, err := sonic.Marshal(args)
	if err != nil {
		return err
	}
	body, err := client.NewBodyFromBytes(jsonBytes)
	if err != nil {
		return err
	}
	req.SetBody(body)

	resp := &client.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return err
	}
	defer resp.Body().Close()
	if resp.IsFail() {
		return resp.ServiceError()
	}
	return nil
}

func RevokeRolePrivileges(cli client.Client, args *RevokeRolePrivilegesArgs) error {
	req := &client.BceRequest{}
	req.SetURI(getRoleURI())
	req.SetMethod(http.Post)
	req.SetParam("revokePrivileges", "")
	jsonBytes, err := sonic.Marshal(args)
	if err != nil {
		return err
	}
	body, err := client.NewBodyFromBytes(jsonBytes)
	if err != nil {
		return err
	}
	req.SetBody(body)

	resp := &client.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return err
	}
	defer resp.Body().Close()
	if resp.IsFail() {
		return resp.ServiceError()
	}
	return nil
}

func ShowRolePrivileges(cli client.Client, args *ShowRolePrivilegesArgs) (*ShowRolePrivilegesResult, error) {
	req := &client.BceRequest{}
	req.SetURI(getRoleURI())
	req.SetMethod(http.Post)
	req.SetParam("showPrivileges", "")
	jsonBytes, err := sonic.Marshal(args)
	if err != nil {
		return nil, err
	}
	body, err := client.NewBodyFromBytes(jsonBytes)
	if err != nil {
		return nil, err
	}
	req.SetBody(body)

	resp := &client.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return nil, err
	}
	defer resp.Body().Close()
	if resp.IsFail() {
		return nil, resp.ServiceError()
	}
	result := &ShowRolePrivilegesResult{}
	if err := resp.ParseJSONBody(result); err != nil {
		return nil, err
	}
	return result, nil
}

func SelectRole(cli client.Client, args *SelectRoleArgs) (*SelectRoleResult, error) {
	req := &client.BceRequest{}
	req.SetURI(getRoleURI())
	req.SetMethod(http.Post)
	req.SetParam("select", "")
	jsonBytes, err := sonic.Marshal(args)
	if err != nil {
		return nil, err
	}
	body, err := client.NewBodyFromBytes(jsonBytes)
	if err != nil {
		return nil, err
	}
	req.SetBody(body)

	resp := &client.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return nil, err
	}
	defer resp.Body().Close()
	if resp.IsFail() {
		return nil, resp.ServiceError()
	}
	result := &SelectRoleResult{}
	if err := resp.ParseJSONBody(result); err != nil {
		return nil, err
	}
	return result, nil
}
