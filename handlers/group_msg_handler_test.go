// @Author Bing 
// @Date 2023/1/17 14:35:00 
// @Desc
package handlers

import (
	"github.com/eatmoreapple/openwechat"
	"github.com/qingconglaixueit/wechatbot/service"
	"reflect"
	"testing"
)

func TestGroupMessageHandler_ReplyText(t *testing.T) {
	type fields struct {
		self    *openwechat.Self
		group   *openwechat.Group
		msg     *openwechat.Message
		sender  *openwechat.User
		service service.UserServiceInterface
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GroupMessageHandler{
				self:    tt.fields.self,
				group:   tt.fields.group,
				msg:     tt.fields.msg,
				sender:  tt.fields.sender,
				service: tt.fields.service,
			}
			if err := g.ReplyText(); (err != nil) != tt.wantErr {
				t.Errorf("ReplyText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGroupMessageHandler_handle(t *testing.T) {
	type fields struct {
		self    *openwechat.Self
		group   *openwechat.Group
		msg     *openwechat.Message
		sender  *openwechat.User
		service service.UserServiceInterface
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GroupMessageHandler{
				self:    tt.fields.self,
				group:   tt.fields.group,
				msg:     tt.fields.msg,
				sender:  tt.fields.sender,
				service: tt.fields.service,
			}
			if err := g.handle(); (err != nil) != tt.wantErr {
				t.Errorf("handle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewGroupMessageHandler(t *testing.T) {
	type args struct {
		msg *openwechat.Message
	}
	tests := []struct {
		name    string
		args    args
		want    MessageHandlerInterface
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGroupMessageHandler(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGroupMessageHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGroupMessageHandler() got = %v, want %v", got, tt.want)
			}
		})
	}
}
