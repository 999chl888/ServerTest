// Copyright 2022 Shopee Inc. All Rights Reserved.
// file_util.go
//
// modification history
// --------------------
// 2022/08/12, created by liruihao
//
// DESCRIPTION:
//     common function related to file

package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Check whether file or directory exists or not
// Args:
//    file_path: path of file or directory
// Returns:
//    true | false, true stands for exists, false stands for not exists
func FileExists(file_path string) bool {
	_, err := os.Stat(file_path)
	return err == nil || os.IsExist(err)
}

// Get file size
// Args:
//    file_path: path of file
// Returns:
//    size: file size in Bytes
//    err: nil stands for success, other stands for not exists
func GetFileSize(file_path string) (int64, error) {
	if !FileExists(file_path) {
		return 0, errors.New(fmt.Sprintf("File(%s) not exists.", file_path))
	}
	file_info, err := os.Stat(file_path)
	if err != nil {
		return 0, err
	}
	return file_info.Size(), nil
}

// check whether file is symbol link
// Args:
//    filepath: file path
// Returns:
//    is_symbol_link: true or false
//    err           : nil stands for success, otherwise stands for fail
func IsSymbolLink(file_path string) (is_symbol_link bool, err error) {
	file_info, err := os.Lstat(file_path)
	if err != nil {
		err_msg := fmt.Sprintf("Fail to get file info of %s: %s", file_path, err.Error())
		return false, errors.New(err_msg)
	}
	return file_info.Mode()&os.ModeSymlink != 0, nil
}

// check path is directory or not
// Args:
//    path: absolute path
// Return:
//    true : path is a directory
//    false: path isn't exist or path isn't a directory
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// Get file's modify time
// Args:
//     file_path: file path
// Returns:
//     modify_time: file's modify time
//     err        : nil stands for success, otherwise stands for fail
func GetFileModifyTime(file_path string) (modify_time time.Time, err error) {
	file_info, err := os.Stat(file_path)
	if err != nil {
		return modify_time, err
	}
	modify_time = file_info.ModTime()
	return modify_time, nil
}

// Get all files in specified directory (contains files in sub directory)
// Args:
//     dir_path: path of directory
// Returns:
//     file_path_list: list of file path
//     err           : nil stands for success, otherwise stands for fail
func GetFilesInDir(dir_path string) (file_path_list []string, err error) {
	if !FileExists(dir_path) {
		err_msg := fmt.Sprintf("%s doesn't exists", dir_path)
		return file_path_list, errors.New(err_msg)
	}

	err = filepath.Walk(dir_path, func(file_path string, file_info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if file_info.IsDir() {
			return nil
		}

		abs_file_path, _ := filepath.Abs(file_path)
		file_path_list = append(file_path_list, abs_file_path)
		return nil
	})
	return file_path_list, err
}
