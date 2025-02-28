package main

import (
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"mime"
	"slices"
	"strings"
	"sync"
	"text/template"

	"github.com/mmcdole/gofeed"
	"github.com/pkierski/pop3srv"
	"golang.org/x/sync/errgroup"
)

type RssMboxProvider struct{}

func (p RssMboxProvider) Provide(user string) (pop3srv.Mailbox, error) {
	e := strings.Split(user, "|")
	slices.Sort(e)
	e = slices.Compact(e)

	if len(e) > 0 && e[0] == "" {
		e = e[:1]
	}

	mbox := MemTableMbox{}
	var mboxMu sync.Mutex

	pool := errgroup.Group{}
	fp := gofeed.NewParser()
	for _, url := range e {
		pool.Go(func() error {
			url = strings.TrimSpace(url)
			feed, err := fp.ParseURL(url)
			if err != nil {
				return nil
			}

			mboxMu.Lock()
			defer mboxMu.Unlock()
			for _, item := range feed.Items {
				mbox.AddWithUidl(makeMail(feed, item), sha256sum(item.GUID))
			}
			return nil
		})
	}

	_ = pool.Wait()

	return &mbox, nil
}

type itemInFeed struct {
	Item *gofeed.Item
	Feed *gofeed.Feed
}

var (
	//go:embed email_item.tmpl
	emailItemTmplStr string
	tmpl             = func() *template.Template {
		tmpl := template.New("item")
		tmpl.Funcs(template.FuncMap{
			"sha256":   sha256sum,
			"qpencode": qpEncodeHeader,
		})
		return template.Must(tmpl.Parse(emailItemTmplStr))
	}()
)

func sha256sum(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func qpEncodeHeader(s string) string {
	return strings.ReplaceAll(
		mime.QEncoding.Encode("utf-8", s),
		" ",
		"\n    ",
	)
}

func makeMail(feed *gofeed.Feed, item *gofeed.Item) string {
	rv := &strings.Builder{}
	tmpl.Execute(rv, itemInFeed{Item: item, Feed: feed})
	return rv.String()
}
