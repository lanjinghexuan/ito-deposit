// Code generated by protoc-gen-go-errors. DO NOT EDIT.

package v1

import (
	fmt "fmt"
	errors "github.com/go-kratos/kratos/v2/errors"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
const _ = errors.SupportPackageIsVersion1

// 未知错误
func IsNearbyUnknownError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == NearbyErrorReason_NEARBY_UNKNOWN_ERROR.String() && e.Code == 500
}

// 未知错误
func ErrorNearbyUnknownError(format string, args ...interface{}) *errors.Error {
	return errors.New(500, NearbyErrorReason_NEARBY_UNKNOWN_ERROR.String(), fmt.Sprintf(format, args...))
}

// 请求参数错误
func IsNearbyBadRequest(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == NearbyErrorReason_NEARBY_BAD_REQUEST.String() && e.Code == 400
}

// 请求参数错误
func ErrorNearbyBadRequest(format string, args ...interface{}) *errors.Error {
	return errors.New(400, NearbyErrorReason_NEARBY_BAD_REQUEST.String(), fmt.Sprintf(format, args...))
}

// 内部服务错误
func IsNearbyInternalError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == NearbyErrorReason_NEARBY_INTERNAL_ERROR.String() && e.Code == 500
}

// 内部服务错误
func ErrorNearbyInternalError(format string, args ...interface{}) *errors.Error {
	return errors.New(500, NearbyErrorReason_NEARBY_INTERNAL_ERROR.String(), fmt.Sprintf(format, args...))
}

// 资源不存在
func IsNearbyNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == NearbyErrorReason_NEARBY_NOT_FOUND.String() && e.Code == 404
}

// 资源不存在
func ErrorNearbyNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(404, NearbyErrorReason_NEARBY_NOT_FOUND.String(), fmt.Sprintf(format, args...))
}
