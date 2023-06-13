package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"net/mail"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/designsbysm/timber/v2"
	"github.com/korylprince/mbox"
)

func processFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	s := mbox.NewScanner(f)
	for s.Scan() {
		b := s.Bytes()

		// copy bytes to buffer, otherwise they will be overwritten
		buf := make([]byte, len(b))
		copy(buf, b)

		// read in mbox separator line
		r := bufio.NewReader(bytes.NewReader(b))
		separator, _, err := r.ReadLine()
		if err != nil {
			return err
		}

		// get entire file
		eml, err := io.ReadAll(r)
		if err != nil {
			return err
		}

		// new reader to headers
		r1 := bytes.NewReader(eml)
		msg, err := mail.ReadMessage(r1)
		if err != nil {
			return err
		}

		// create the folder
		date, err := msg.Header.Date()
		if err != nil {
			// if no date header, use the separator date as a fallback
			regex, _ := regexp.Compile(`.{3} .{3} \d{2} \d{2}:\d{2}:\d{2} \d{4}$`)
			match := regex.Find(separator)

			if len(match) > 0 {
				date, err = time.Parse(time.ANSIC, string(match))
				if err != nil {
					return err
				}
			}
		}
		folder := fmt.Sprintf(date.Format("%s/2006/01 - January"), filepath.Dir(filename))
		if err := os.MkdirAll(folder, 0755); err != nil {
			return err
		}

		// save the file
		file := fmt.Sprintf(date.Format("%s/2006-01-02 15:04:05 %s.eml"), folder, msg.Header.Get("Subject"))
		timber.Info("Saving:", file)
		os.WriteFile(file, eml, fs.ModePerm)
	}

	return nil
}
