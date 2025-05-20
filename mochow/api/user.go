package api

import (
	"github.com/baidu/mochow-sdk-go/v2/client"
	"github.com/baidu/mochow-sdk-go/v2/http"
	"github.com/bytedance/sonic"
)

func CreateUser(cli client.Client, args *CreateUserArgs) error {
	req := &client.BceRequest{}
	req.SetURI(getUserURI())
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

func DropUser(cli client.Client, args *DropUserArgs) error {
	req := &client.BceRequest{}
	req.SetURI(getUserURI())
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

func ChangeUserPassword(cli client.Client, args *ChangeUserPasswordArgs) error {
	req := &client.BceRequest{}
	req.SetURI(getUserURI())
	req.SetMethod(http.Post)
	req.SetParam("changePassword", "")
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

func GrantUserRoles(cli client.Client, args *GrantUserRolesArgs) error {
	req := &client.BceRequest{}
	req.SetURI(getUserURI())
	req.SetMethod(http.Post)
	req.SetParam("grantRoles", "")
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

func RevokeUserRoles(cli client.Client, args *RevokeUserRolesArgs) error {
	req := &client.BceRequest{}
	req.SetURI(getUserURI())
	req.SetMethod(http.Post)
	req.SetParam("revokeRoles", "")
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

func GrantUserPrivileges(cli client.Client, args *GrantUserPrivilegesArgs) error {
	req := &client.BceRequest{}
	req.SetURI(getUserURI())
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

func RevokeUserPrivileges(cli client.Client, args *RevokeUserPrivilegesArgs) error {
	req := &client.BceRequest{}
	req.SetURI(getUserURI())
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

func ShowUserPrivileges(cli client.Client, args *ShowUserPrivilegesArgs) (*ShowUserPrivilegesResult, error) {
	req := &client.BceRequest{}
	req.SetURI(getUserURI())
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
	result := &ShowUserPrivilegesResult{}
	if err := resp.ParseJSONBody(result); err != nil {
		return nil, err
	}
	return result, nil
}

func SelectUser(cli client.Client, args *SelectUserArgs) (*SelectUserResult, error) {
	req := &client.BceRequest{}
	req.SetURI(getUserURI())
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
	result := &SelectUserResult{}
	if err := resp.ParseJSONBody(result); err != nil {
		return nil, err
	}
	return result, nil
}
