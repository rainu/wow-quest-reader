package system

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

type teeReader struct {
	source     io.ReadCloser
	output     io.WriteCloser
	outputPath string
}

func NewTeeReader(source io.ReadCloser, outPath string) (io.ReadCloser, error) {
	f, err := os.OpenFile(outPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}

	return &teeReader{
		source:     source,
		output:     f,
		outputPath: outPath,
	}, nil
}

func (t *teeReader) Read(p []byte) (n int, err error) {
	n, err = t.source.Read(p)

	_, wErr := t.output.Write(p[0:n])
	if wErr != nil {
		logrus.WithField("file", t.outputPath).WithError(err).Warnf("Unable to write stream to output file!")
	}

	return n, err
}

func (t *teeReader) Close() error {
	defer t.output.Close()

	// copy the potential rest of stream to file!
	io.Copy(t.output, t.source)

	return t.source.Close()
}
