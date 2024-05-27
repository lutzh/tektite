package opers

import (
	"github.com/spirit-labs/tektite/encoding"
	"github.com/spirit-labs/tektite/evbatch"
	"github.com/spirit-labs/tektite/types"
	"sync"
)

type StoreStreamOperator struct {
	BaseOperator
	inSchema      *OperatorSchema
	outSchema     *OperatorSchema
	keyCols       []int
	rowCols       []int
	slabID        int
	keyPrefix     []byte
	nodeID        int
	addOffset     bool
	nextOffsets   []int64
	store         store
	offsetsSlabID int
}

func NewStoreStreamOperator(inSchema *OperatorSchema, slabID int, offsetsSlabID int, store store, nodeID int) (*StoreStreamOperator, error) {
	nextOffsets := make([]int64, inSchema.PartitionScheme.Partitions)
	for i := range nextOffsets {
		nextOffsets[i] = -1
	}
	var outEventSchema *evbatch.EventSchema
	addOffset := offsetsSlabID != -1
	if addOffset {
		colTypes := []types.ColumnType{types.ColumnTypeInt}
		colTypes = append(colTypes, inSchema.EventSchema.ColumnTypes()...)
		colNames := []string{OffsetColName}
		colNames = append(colNames, inSchema.EventSchema.ColumnNames()...)
		outEventSchema = evbatch.NewEventSchema(colNames, colTypes)
	} else {
		outEventSchema = inSchema.EventSchema
	}
	var rowCols []int
	for i := 1; i < len(outEventSchema.ColumnNames()); i++ {
		rowCols = append(rowCols, i)
	}
	outSchema := inSchema.Copy()
	outSchema.EventSchema = outEventSchema

	keyPrefix := encoding.AppendUint64ToBufferBE(nil, uint64(slabID))
	return &StoreStreamOperator{
		inSchema:      inSchema,
		outSchema:     outSchema,
		slabID:        slabID,
		keyPrefix:     keyPrefix,
		keyCols:       []int{0},
		rowCols:       rowCols,
		nodeID:        nodeID,
		addOffset:     addOffset,
		store:         store,
		offsetsSlabID: offsetsSlabID,
		nextOffsets:   nextOffsets,
	}, nil
}

func (ts *StoreStreamOperator) HandleQueryBatch(*evbatch.Batch, QueryExecContext) (*evbatch.Batch, error) {
	panic("not supported in queries")
}

func (ts *StoreStreamOperator) HandleStreamBatch(batch *evbatch.Batch, execCtx StreamExecContext) (*evbatch.Batch, error) {
	if ts.addOffset {
		// Add the offset column
		colBuilder := evbatch.NewIntColBuilder()
		kOffset, err := ts.getNextOffset(execCtx.PartitionID())
		if err != nil {
			return nil, err
		}
		rc := batch.RowCount
		for i := 0; i < rc; i++ {
			colBuilder.Append(kOffset)
			kOffset++
		}
		newOffsetCol := colBuilder.Build()
		for _, col := range batch.Columns {
			col.Retain()
		}
		newCols := make([]evbatch.Column, len(batch.Columns)+1)
		newCols[0] = newOffsetCol
		copy(newCols[1:], batch.Columns)
		batch = evbatch.NewBatch(ts.outSchema.EventSchema, newCols...)

		ts.nextOffsets[execCtx.PartitionID()] = kOffset
		storeOffset(execCtx, kOffset, ts.offsetsSlabID, execCtx.WriteVersion())
	}

	keyPrefix := encoding.AppendUint64ToBufferBE(ts.keyPrefix, uint64(execCtx.PartitionID()))
	storeBatchInTable(batch, []int{0}, ts.rowCols, keyPrefix, execCtx, ts.nodeID, false)
	return batch, ts.sendBatchDownStream(batch, execCtx)
}

func (ts *StoreStreamOperator) HandleBarrier(execCtx StreamExecContext) error {
	return ts.BaseOperator.HandleBarrier(execCtx)
}

func (ts *StoreStreamOperator) getNextOffset(partitionID int) (int64, error) {
	nextOffset := ts.nextOffsets[partitionID]
	if nextOffset != -1 {
		return nextOffset, nil
	}
	// Load from store
	return loadOffset(ts.offsetsSlabID, partitionID, ts.store)
}

func (ts *StoreStreamOperator) InSchema() *OperatorSchema {
	return ts.inSchema
}

func (ts *StoreStreamOperator) OutSchema() *OperatorSchema {
	return ts.outSchema
}

func (ts *StoreStreamOperator) Setup(StreamManagerCtx) error {
	return nil
}

func (ts *StoreStreamOperator) Teardown(StreamManagerCtx, *sync.RWMutex) {
}
