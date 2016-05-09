package memory

import (
	"github.com/emersion/imap/common"
)

type Mailbox struct {
	name string
	messages []*Message
}

func (mbox *Mailbox) Info() (*common.MailboxInfo, error) {
	info := &common.MailboxInfo{
		Delimiter: "/",
		Name: mbox.name,
	}
	return info, nil
}

func (mbox *Mailbox) uidNext() (uid uint32) {
	for _, msg := range mbox.messages {
		if msg.Uid > uid {
			uid = msg.Uid
		}
	}
	uid++
	return
}

func (mbox *Mailbox) Status(items []string) (*common.MailboxStatus, error) {
	status := &common.MailboxStatus{
		Items: items,
		Name: mbox.name,
	}

	for _, name := range items {
		switch name {
		case "MESSAGES":
			status.Messages = uint32(len(mbox.messages))
		case "UIDNEXT":
			status.UidNext = mbox.uidNext()
		}
	}

	return status, nil
}

func (mbox *Mailbox) ListMessages(uid bool, seqset *common.SeqSet, items []string) (msgs []*common.Message, err error) {
	for i, msg := range mbox.messages {
		id := uint32(i+1)

		if (uid && !seqset.Contains(msg.Uid)) || (!uid && !seqset.Contains(id)) {
			continue
		}

		m := msg.Metadata(items)
		m.Id = id
		msgs = append(msgs, m)
	}

	return
}

func (mbox *Mailbox) SearchMessages(uid bool, criteria *common.SearchCriteria) (ids []uint32, err error) {
	for i, msg := range mbox.messages {
		if !msg.Matches(criteria) {
			continue
		}

		var id uint32
		if uid {
			id = msg.Uid
		} else {
			id = uint32(i+1)
		}
		ids = append(ids, id)
	}

	return
}