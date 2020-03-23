
package semtools


import (
    "github.com/sirupsen/logrus"
)

func init() {

    logrus.SetFormatter(&logrus.TextFormatter{
        DisableColors: true,
        FullTimestamp: true,
    })

}


func GetLogger(context string) *logrus.Entry {
    return logrus.WithFields(logrus.Fields{"context": context})
}