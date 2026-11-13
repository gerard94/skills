/* 
Skills

Copyright (C) 2023 GérardMeunier

This program is free software; you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation; either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License  for more details.

You should have received a copy of the GNU General Public License along with this program; if not, write to the Free Software Foundation, Inc., 59 Temple Place - Suite 330, Boston, MA  02111-1307, USA.
*/

package gqReceiver

import (
	
	A	"github.com/gerard94/util/avl"
	B	"github.com/gerard94/skills/base"
	BA	"github.com/gerard94/skills/basic"
	F	"path/filepath"
	G	"github.com/gerard94/util/graphQL"
	J	"github.com/gerard94/util/json"
	M	"github.com/gerard94/util/misc"
	Q	"github.com/gerard94/util/queue"
	SC	"strconv"
		
		"errors"
		"fmt"
		"net/http"
		"io/ioutil"
		"os"
		"os/signal"
		"strings"
		"sync"
		"syscall"
		"time"

)

const (
	
	typeSystemName = "typeSystem.txt"
	
	activateDelay = 200 * time.Millisecond

)

type (
	
	action struct {
		es G.ExecSystem
		doc *G.Document
		opName string
		varVals J.Json
		variableValues *A.Tree // *ValMapItem
		w http.ResponseWriter
		c chan bool
	}
	
	scalarer struct {
	}
	
	typeSystem struct {
		G.TypeSystem
	}

)

var (
	
	lg = M.GetLog()

)

var (
	
	sl scalarer
	
	ts *typeSystem
	
	typeSystemPath = F.Join(BA.GetSystemPath(), typeSystemName)
	
	serverAddress = BA.ServerAddress()
	
	opMutex sync.RWMutex

)

func printErrors (errors *A.Tree) {
	for e := errors.Next(nil); e != nil; e = errors.Next(e) {
		el := e.Val().(*G.ErrorElem)
		var b strings.Builder
		b.WriteString("***ERROR*** " + el.Message)
		if el.Locations != nil {
			b.WriteString(" at")
			for i, loc := range el.Locations {
				if i > 0 {
					b.WriteString(",")
				}
				b.WriteString(" ")
				b.WriteString(SC.Itoa(loc.L))
				b.WriteString(":")
				b.WriteString(SC.Itoa(loc.C))
			}
		}
		p := el.Path
		if p != nil {
			b.WriteString(" in ")
			for {
				switch p := p.(type) {
				case *G.PathString:
					b.WriteString(p.S)
				case *G.PathNb:
					b.WriteString(SC.Itoa(p.N))
				}
				p = p.Next()
				if p == nil {
					break
				}
				b.WriteString(".")
			}
		}
		s := b.String()
		lg.Println(s)
		/*
		fmt.Println(s)
		*/
	}
} //printErrors

func (a *action) activate () {
	a.w.Header().Set("Content-Type", "application/json")
	//a.w.WriteHeader(http.StatusCreated)
	var r G.Response
	switch a.es.GetOperation(a.opName).OpType {
	case G.QueryOp:
		opMutex.RLock()
		r = a.es.Execute(a.doc, a.opName, a.variableValues)
		opMutex.RUnlock()
	case G.MutationOp:
		opMutex.Lock()
		r = a.es.Execute(a.doc, a.opName, a.variableValues)
		opMutex.Unlock()
	default:
		M.Halt(100)
	}
	errors := r.Errors()
	if errors != nil {
		printErrors(errors)
	}
	G.ResponseToJson(r).Write(a.w)
	a.c <- true
} //Activate

func (a *action) name () string {
	if a.opName == "" {
		return "anonymous"
	} else {
		return a.opName
	}
} //Name

func readOpNameVars  (req *http.Request) (varVals J.Json, t *A.Tree, opName, docS string, err error) {
	buf, error := ioutil.ReadAll(req.Body); M.Assert(error == nil, error, 100)
	sB := string(buf)
	j := J.ReadString(sB)
	if j == nil || j.JKind() != J.ObjectKind {
		s := "Incorrect JSON request: "
		if sB != "" {
			s += sB
		} else {
			s += "Void string"
		}
		err = errors.New(s)
		return
	}
	o := j.(*J.Object)
	v, _ := J.GetValue(o, "operationName")
	opName, _ = J.GetString(v)
	v, b := J.GetValue(o, "variables")
	if b {
		varVals, b = J.GetJson(v)
	}
	if !b {
		varVals = J.ReadString("{}")
	}
	if j.JKind() != J.ObjectKind {
		err = errors.New("Incorrect variables value")
		return
	}
	obj := varVals.(*J.Object)
	t = A.New()
	for _, f := range obj.Fields {
		if !G.InsertJsonValue(t, f.Name, f.Value) {
			err = errors.New("Duplicated variable name " + f.Name)
			return
		}
	}
	v, _ = J.GetValue(o, "query")
	docS, _ = J.GetString(v)
	if docS == "" {
		err = errors.New("No query string")
	}
	return
} //readOpNameVars

func makeHandler (newAction chan<- *action) func (w http.ResponseWriter, req *http.Request) {
	
	return func (w http.ResponseWriter, req *http.Request) {
		
		writeError := func (err error) {
			m := J.NewMaker()
			m.StartObject()
			m.PushString(err.Error())
			m.BuildField("errors")
			m.BuildObject()
			m.GetJson().Write(w)
			lg.Println("***ERROR*** ", err)
		}
		
		if req.Method != "GET" && req.Method != "POST" {
			return
		}
		
		j, variableValues, opName, docS, error := readOpNameVars (req)
		if error != nil {
			writeError(error)
			return
		}
		doc, r := G.ReadString(docS)
		if doc == nil {
			err := r.Errors()
			M.Assert(!err.IsEmpty(), 102)
			printErrors(err)
			G.ResponseToJson(r).Write(w)
			return
		}
		if !G.ExecutableDefinitions(doc) {
			ts.MappedError("#skills:NotExecDefs", "", "", nil, nil)
			err := ts.GetErrors()
			printErrors(err)
			r := new(G.InstantResponse)
			r.SetErrors(err)
			G.ResponseToJson(r).Write(w)
			return
		}
		es := ts.ExecValidate(doc)
		err := es.GetErrors()
		if !err.IsEmpty() {
			printErrors(err)
			r := new(G.InstantResponse)
			r.SetErrors(err)
			G.ResponseToJson(r).Write(w)
			return
		}
		if opName == "" {
			if opList := es.ListOperations(); len(opList) == 1 {
				opName = opList[0]
			} else {
				writeError(errors.New("Selected operation name not defined"))
				return
			}
		}
		a := &action{es: es, doc: doc, opName: opName, varVals: j, variableValues: variableValues, w: w, c: make(chan bool)}
		newAction <- a
		<- a.c
	}

} //makeHandler

func loop (newAction chan<- *action) {
	r := http.NewServeMux()
	r.HandleFunc("/", makeHandler(newAction))
	server := &http.Server{
		Addr: serverAddress,
		Handler: r,
	}
	s := fmt.Sprint("Listening on ", serverAddress, " ...")
	lg.Println(s)
	fmt.Println(s)
	server.ListenAndServe()
} //loop

func coercePubkey (ts G.TypeSystem, v G.Value, path G.Path, cV *G.Value) bool {
	switch v := v.(type) {
	case *G.StringValue:
		*cV = v
		return true
	default:
		ts.MappedError("NotAPubkey", "", "", nil, path)
		return false
	}
} //coercePubkey

func (scalarer) FixScalarCoercer (scalarName string, sc *G.ScalarCoercer) {
	switch scalarName {
	case "Pubkey":
		*sc = coercePubkey
	}
} //FixScalarCoercer

var q = Q.NewQueue[*action]()

func manageActions (newAction <-chan *action) {
	for {
		select {
		case a := <-newAction:
			q.Push(a)
		default:
			for !q.IsEmpty() {
				a := q.Pop()
				a.activate()
			}
		}
		time.Sleep(activateDelay)
	}
} //manageActions

func Start () {
	newAction := make(chan *action)
	lg.Println("Starting"); lg.Println()
	B.OpenBase()
	defer B.CloseBase()
	stopProg := make(chan os.Signal, 1)
	signal.Notify(stopProg, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	go loop(newAction)
	go manageActions(newAction)
	<- stopProg
	opMutex.Lock()
	lg.Println("Halting"); lg.Println()
} //Start

func TS () G.TypeSystem {
	return ts
} //TS


func init () {

	const (
		
		rootName = "root"

	)
	
	doc, r := G.ReadGraphQL(typeSystemPath)
	err := r.Errors()
	if !err.IsEmpty() {
		printErrors(err)
		M.Halt(100)
	}
	M.Assert(doc != nil, 101)
	if G.ExecutableDefinitions(doc) {
		ts.MappedError("#duniter:NoTypeSystemInDoc", "", "", nil, nil)
		err := ts.GetErrors()
		printErrors(err)
		M.Halt(102)
	}
	tsRead := false
	ts = new(typeSystem)
	ts.TypeSystem = G.Dir.NewTypeSystem(sl)
	ts.InitTypeSystem(doc)
	initialValue := G.NewOutputObjectValue()
	initialValue.InsertOutputField(rootName, nil)
	ts.FixInitialValue(initialValue)
	tsRead = ts.GetErrors().IsEmpty()
	if !tsRead {
		err := ts.GetErrors()
		printErrors(err)
		M.Halt(103)
	}
} //init
