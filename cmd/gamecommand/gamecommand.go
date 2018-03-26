package gamecommand

import (
	"bufio"
	"os"
	"os/signal"
	"syscall"

	"github.com/UnnecessaryRain/ironway-core/pkg/game"
	"github.com/UnnecessaryRain/ironway-core/pkg/game/commands"
	log "github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type gameCommand struct {
	player string
}

// Configure sets up the command for gamer
func Configure(app *kingpin.Application) {
	g := &gameCommand{}
	c := app.Command("game", "starts a game").
		Action(g.run)
	c.Flag("player", "player to use for this game").
		Short('p').
		Default("irony42").
		StringVar(&g.player)
}

func (g *gameCommand) run(c *kingpin.ParseContext) error {
	log.Infoln("Starting standalone game with player", g.player)

	// signal handling to shutdown gracefully
	sigs := make(chan os.Signal)
	stopChan := make(chan struct{})
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		close(stopChan)
		// TODO: close http gamer gracefully aswell
		os.Exit(0)
	}()

	gameInstance := game.NewGame()
	go gameInstance.RunForever(stopChan)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cmd := commands.NewDebug(scanner.Text())
		gameInstance.QueueCommand(cmd)
	}

	close(stopChan)

	return nil
}