package indexer

import (
	"fmt"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/celo-indexer/store/psql"

	"github.com/figment-networks/celo-indexer/model"
	"github.com/pkg/errors"
)

// reportCreator creates and completes report
type reportCreator struct {
	kind         model.ReportKind
	indexVersion int64
	startHeight  int64
	endHeight    int64

	reportDb store.Reports

	report *model.Report
}

func (o *reportCreator) createIfNotExists(kinds ...model.ReportKind) error {
	report, err := o.reportDb.FindNotCompletedByIndexVersion(o.indexVersion, kinds...)
	if err != nil {
		if err == psql.ErrNotFound {
			if err = o.create(); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		if report.Kind != o.kind {
			return errors.New(fmt.Sprintf("there is already reindexing in process [kind=%s] (use -force flag to override it)", report.Kind))
		}
		o.report = report
	}
	return nil
}

func (o *reportCreator) create() error {
	report := &model.Report{
		Kind:         o.kind,
		IndexVersion: o.indexVersion,
		StartHeight:  o.startHeight,
		EndHeight:    o.endHeight,
	}

	if err := o.reportDb.Create(report); err != nil {
		return err
	}

	o.report = report

	return nil
}

func (o *reportCreator) complete(totalCount int64, successCount int64, err error) error {
	o.report.Complete(successCount, totalCount-successCount, err)

	return o.reportDb.Save(o.report)
}
