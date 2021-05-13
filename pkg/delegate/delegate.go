package delegate

type Delegate struct {
	ObjectCreator CreateObjectFunc
	ListCreator   CreateObjectFunc
	DBProvider    DBFunc
	PKExtractor   ExtractPKFunc
	PKVerificator VerifyPKFunc
	PKAssigner    AssignPKFunc
}
