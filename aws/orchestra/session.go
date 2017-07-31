package orchestra

import (
	//"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
)

func CreateSession() (*session.Session, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	if _, err := sess.Config.Credentials.Get(); err != nil {
		return nil, err
	}
	return sess, nil
}
