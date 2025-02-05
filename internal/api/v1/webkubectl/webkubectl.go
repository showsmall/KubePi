package webkubectl

import (
	"fmt"
	"github.com/KubeOperator/kubepi/internal/api/v1/session"
	"github.com/KubeOperator/kubepi/internal/service/v1/cluster"
	"github.com/KubeOperator/kubepi/internal/service/v1/clusterbinding"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/KubeOperator/kubepi/pkg/kubernetes"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type Handler struct {
	clusterBindingService clusterbinding.Service
	clusterService        cluster.Service
	sessionCache          *TerminalSessions
}

func NewHandler() *Handler {
	return &Handler{
		clusterBindingService: clusterbinding.NewService(),
		clusterService:        cluster.NewService(),
		sessionCache:          NewTerminalSessions(),
	}
}

func (h *Handler) GetConfigFile() iris.Handler {
	return func(ctx *context.Context) {
		sessionId := ctx.URLParam("token")
		requireLen := len(uuid.New().String())
		if len(sessionId) != requireLen {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Sprintf("sessionId length must be %d", requireLen))
			return
		}
		sess := h.sessionCache.Get(sessionId)
		if sess != nil {
			h.sessionCache.Delete(sessionId)
		} else {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("can not find sessionId: %s in memory", sessionId))
			return
		}

		ctx.Header("Content-Type", "application/download")
		ctx.Header("Content-Disposition", "attachment;filename=config")
		ctx.Header("Content-Transfer-Encoding", "binary")

		cc := toCmdConfig(sess)
		bs, err := clientcmd.Write(*cc)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", "can not generate config file")
			return
		}
		_, _ = ctx.Write(bs)
	}
}

func toCmdConfig(sess *Session) *clientcmdapi.Config {
	cc := clientcmdapi.NewConfig()
	cc.Clusters[sess.Cluster] = &clientcmdapi.Cluster{
		Server:                sess.config.Host,
		InsecureSkipTLSVerify: true,
	}
	cc.AuthInfos[sess.User] = &clientcmdapi.AuthInfo{
		ClientCertificateData: sess.config.CertData,
		ClientKeyData:         sess.config.KeyData,
		Token:                 sess.config.BearerToken,
	}
	contextName := fmt.Sprintf("%s@%s", sess.Cluster, sess.User)
	cc.Contexts[contextName] = &clientcmdapi.Context{
		Cluster:  sess.Cluster,
		AuthInfo: sess.User,
	}
	cc.CurrentContext = contextName
	return cc
}

func (h *Handler) CreateSession() iris.Handler {
	return func(ctx *context.Context) {
		var sess Session
		if err := ctx.ReadJSON(&sess); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		u := ctx.Values().Get("profile")
		profile := u.(session.UserProfile)

		c, err := h.clusterService.Get(sess.Cluster, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		k := kubernetes.NewKubernetes(c)
		cfg, err := k.Config()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		if !profile.IsAdministrator {
			rb, err := h.clusterBindingService.GetBindingByClusterNameAndUserName(sess.Cluster, profile.Name, common.DBOptions{})
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
			cfg.CertData = rb.Certificate
		}
		sess.config = cfg
		sess.User = profile.Name
		sessionId := uuid.New().String()
		h.sessionCache.Put(sessionId, &sess)
		ctx.Values().Set("data", &SessionResponse{Token: sessionId})
	}
}

func Install(authParent, noAuthParty iris.Party) {
	handler := NewHandler()
	authParent.Post("/webkubectl/session", handler.CreateSession())
	noAuthParty.Get("/webkubectl/session", handler.GetConfigFile())
}
