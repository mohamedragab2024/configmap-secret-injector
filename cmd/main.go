package main

import (
	"configmap-secret-injector/internal/controllers"
	"flag"
	"os"
	"time"

	"github.com/rs/zerolog"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	scheme               = runtime.NewScheme()
	enableLeaderElection bool
	leaderElectionID     string
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

}
func main() {
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.StringVar(&leaderElectionID, "leader-elect-id", "configmap-secret-injector",
		"Leader election ID for controller manager. "+
			"this will ensure there is only one active controller manager.")

	opts := zap.Options{
		DestWriter:  os.Stdout,
		Development: true,
		Encoder: zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			MessageKey:  "msg",
			LevelKey:    "level",
			TimeKey:     "ts",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
			EncodeLevel: zapcore.CapitalLevelEncoder,
		}),
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	logger := zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(zerolog.TraceLevel).With().Timestamp().Caller().Logger()
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:           scheme,
		LeaderElection:   enableLeaderElection,
		LeaderElectionID: leaderElectionID,
	})

	if err != nil {
		logger.Error().Err(err).Msg("unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.ConfigMapReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Logger: &logger,
	}).New(mgr); err != nil {
		logger.Error().Err(err).Msg("unable to create controller")
	}
	logger.Info().Msg("Starting manager")

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		logger.Error().Err(err).Msg("problem running manager")
		os.Exit(1)
	}
}
