/* 
Skills

Copyright (C) 2023 GérardMeunier

This program is free software; you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation; either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License  for more details.

You should have received a copy of the GNU General Public License along with this program; if not, write to the Free Software Foundation, Inc., 59 Temple Place - Suite 330, Boston, MA  02111-1307, USA.
*/

package basic

// Package basic implements general definitions and utilities for the database server

import (
	
	F	"path/filepath"
	M	"github.com/gerard94/util/misc"
	R	"github.com/gerard94/util/resources"
		"errors"
		"fmt"
		"log"
		"os"
		"strconv"
		"text/scanner"

)

const (
	
	rsrcDir = "skills"
	systemDir = "skills/system"
	
	serverDefaultAddress = "localhost:8081"
	// Name of the file which contains the server address 
	serverAddressName = "serverAddress.txt"
	
	pageNbDefault = 2000
	// Name of the file in the 'rsrcDir' resource directory which contains the highest number of buffer pages in indexes 
	pageNbName = "baseBufferSize.txt"
	
	// Logs; when the current log contains more than 'maxLogSize' bytes, it is renamed as old and a new current log is started
	logName = "log.txt"
	logOldName = "log1.txt"
	maxLogSize = 0x3200000 // 50MB

)

var (
	
	// Paths in the 'rsrc' directory
	rsrcPath = R.FindDir(rsrcDir)
	systemPath = R.FindDir(systemDir)
	
	logPath = F.Join(rsrcPath, logName)
	logOldPath = F.Join(rsrcPath, logOldName)
	
	// Current log
	lg *log.Logger
	
	serverAddress = serverDefaultAddress
	
	// Number of pages in the buffer of the database
	pageNb = pageNbDefault

)

func GetRsrcPath () string {
	return rsrcPath
}

func GetSystemPath () string {
	return systemPath
}

func setLog () {
	fi, err := os.Stat(logPath)
	M.Assert(err == nil || os.IsNotExist(err), err, 100)
	var f *os.File
	// if current log is too big at startup, rename it as old and start a new one
	if err != nil || fi.Size() >= maxLogSize {
		err := os.Remove(logOldPath)
		M.Assert(err == nil || os.IsNotExist(err), err, 101)
		err = os.Rename(logPath, logOldPath)
		M.Assert(err == nil || os.IsNotExist(err), err, 102)
		f, err = os.Create(logPath)
		M.Assert(err == nil, err, 103)
	} else { // Open the current log in append mode
		f, err = os.OpenFile(logPath, os.O_APPEND | os.O_WRONLY, 0644)
		M.Assert(err == nil, err, 104)
	}
	lg = log.New(f, "", log.Ldate | log.Ltime | log.Lshortfile)
	M.SetLog(lg)
} //setLog

func ServerAddress () string {
	return serverAddress
}

// Fix server address at startup
func fixServerAddress () {
	name := F.Join(rsrcPath, serverAddressName)
	f, err := os.Open(name)
	if err == nil {
		defer f.Close()
		s := new(scanner.Scanner)
		s.Init(f)
		s.Error = func(s *scanner.Scanner, msg string) {M.Halt(errors.New("File " + name + " incorrect"), 100)}
		s.Mode = scanner.ScanStrings
		s.Scan()
		ss := s.TokenText()
		M.Assert(ss[0] == '"' && ss[len(ss) - 1] == '"', ss, 100)
		serverAddress = ss[1:len(ss) - 1]
	} else {
		f, err := os.Create(name)
		M.Assert(err == nil, err, 101)
		defer f.Close()
		fmt.Fprint(f, "\"" + serverAddress + "\"")
	}
}

func PageNb () int {
	return pageNb
}

// Fix the highest buffer page number of indexes
func fixPageNb () {
	name := F.Join(rsrcPath, pageNbName)
	f, err := os.Open(name)
	if err == nil {
		defer f.Close()
		s := new(scanner.Scanner)
		s.Init(f)
		s.Error = func(s *scanner.Scanner, msg string) {M.Halt(errors.New("File " + name + " incorrect"), 100)}
		s.Mode = scanner.ScanInts
		s.Scan()
		ss := s.TokenText()
		pageNb, err = strconv.Atoi(ss)
		M.Assert(err == nil, err, 100)
	} else {
		f, err := os.Create(name)
		M.Assert(err == nil, err, 101)
		defer f.Close()
		fmt.Fprint(f, pageNb)
	}
}

func init () {
	err := os.MkdirAll(rsrcPath, 0777); M.Assert(err == nil, err, 100)
	err = os.MkdirAll(systemPath, 0777); M.Assert(err == nil, err, 101)
	setLog()
	fixServerAddress()
	fixPageNb()
} //init
