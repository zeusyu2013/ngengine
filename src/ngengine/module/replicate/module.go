// package replica
// 对象同步模块
package replicate

import (
	"fmt"
	"ngengine/core/service"
	"ngengine/module/object"
)

type ReplicateModule struct {
	core           service.CoreApi
	objectmodule   *object.ObjectModule
	replicate      *Replicate
	replicatetypes map[string]struct{}
	delayregs      []string
}

func New() *ReplicateModule {
	o := &ReplicateModule{}
	o.replicate = NewReplicate(o)
	o.replicatetypes = make(map[string]struct{})
	o.delayregs = make([]string, 0, 16)
	return o
}

// Name 模块名
func (o *ReplicateModule) Name() string {
	return "Replicate"
}

// Init 模块初始化
func (o *ReplicateModule) Init(core service.CoreApi) bool {
	o.core = core
	o.objectmodule = core.Module("Object").(*object.ObjectModule)
	if o.objectmodule == nil {
		panic("need object module")
	}

	return true
}

// Start 模块启动
func (o *ReplicateModule) Start() {
	for _, typ := range o.delayregs {
		err := o.objectmodule.AddEventCallback(typ, "on_create", o.replicate.ObjectCreate, object.PRIORITY_LOWEST)
		if err != nil {
			panic(err)
		}
		err = o.objectmodule.AddEventCallback(typ, "on_destroy", o.replicate.ObjectDestroy, object.PRIORITY_LOWEST)
		if err != nil {
			panic(err)
		}
	}

	o.delayregs = o.delayregs[:0]
}

// Shut 模块关闭
func (o *ReplicateModule) Shut() {

}

// OnUpdate 模块Update
func (o *ReplicateModule) OnUpdate(t *service.Time) {
}

// OnMessage 模块消息
func (o *ReplicateModule) OnMessage(id int, args ...interface{}) {
}

// RegisterReplicate 注册需要同步的对象类型
func (o *ReplicateModule) RegisterReplicate(typ string) error {
	if _, has := o.replicatetypes[typ]; has {
		return fmt.Errorf("register replicate twice, %s", typ)
	}

	o.replicatetypes[typ] = struct{}{}
	if o.objectmodule != nil {
		err := o.objectmodule.AddEventCallback(typ, "on_create", o.replicate.ObjectCreate, object.PRIORITY_LOWEST)
		if err != nil {
			return err
		}
		err = o.objectmodule.AddEventCallback(typ, "on_destroy", o.replicate.ObjectDestroy, object.PRIORITY_LOWEST)
		return err
	}

	o.delayregs = append(o.delayregs, typ)
	return nil
}
