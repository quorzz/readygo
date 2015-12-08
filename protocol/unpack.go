package protocol

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

type ResponseType int

const (
	ResponseError ResponseType = iota
	ResponseNil
	ResponseStatus
	ResponseInt
	ResponseBulk
	ResponseMutli
)

type Response struct {
	Type    ResponseType
	Error   string
	Status  string
	Integer int64
	Bulk    []byte
	Multi   []*Response
}

func ReadResponse(r *bufio.Reader) (*Response, error) {
	line, e := r.ReadBytes('\n')
	if e != nil {
		return nil, e
	}

	line = line[:len(line)-2]
	switch line[0] {
	case '-':
		return &Response{
			Type:  ResponseError,
			Error: string(line[1:]),
		}, nil

	case '+':
		return &Response{
			Type:   ResponseStatus,
			Status: string(line[1:]),
		}, nil

	case ':':
		n, err := strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return nil, err
		}

		return &Response{
			Type:    ResponseInt,
			Integer: n,
		}, nil

	case '$':
		l, err := strconv.Atoi(string(line[1:]))
		if err != nil {
			return nil, err
		}

		if l < 0 {
			return nil, nil
		}
		buf := make([]byte, l+2)
		if _, err := io.ReadFull(r, buf); err != nil {
			return nil, err
		}
		return &Response{
			Bulk: buf[:l],
			Type: ResponseBulk,
		}, nil

	case '*':
		l, e := strconv.Atoi(string(line[1:]))
		if e != nil {
			return nil, e
		}

		if l < 0 {
			return &Response{Type: ResponseMutli}, nil
		}
		ret := make([]*Response, l)
		for i := 0; i < l; i++ {
			m, err := ReadResponse(r)
			if err != nil {
				return nil, err
			}
			ret[i] = m
		}
		return &Response{
			Type:  ResponseMutli,
			Multi: ret,
		}, nil
	}
	return nil, errors.New("redis protocol errors")
}

func (response Response) ToString() (string, error) {
	if response.Type == ResponseError {
		return "", errors.New(response.Error)
	}

	if response.Type != ResponseBulk {
		return "", errors.New("not type of bulk")
	}

	if response.Bulk == nil {
		return "", nil
	}
	return string(response.Bulk), nil
}

func (response Response) ToMap() (map[string]string, error) {
	if response.Type != ResponseMutli {
		return nil, errors.New("not type of multi")
	}

	result := make(map[string]string)
	length := len(response.Multi)
	if response.Multi == nil || length <= 0 {
		return result, nil
	}

	for i := 0; i < length/2; i++ {
		key, err := response.Multi[i*2].ToString()
		if err != nil {
			return nil, err
		}

		value, err := response.Multi[1*2+1].ToString()
		if err != nil {
			return nil, err
		}

		result[key] = value
	}
	return result, nil
}
