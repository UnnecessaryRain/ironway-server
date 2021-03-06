package gamecmd

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/UnnecessaryRain/ironway-core/pkg/mud/game"
	"github.com/UnnecessaryRain/ironway-core/pkg/mud/interpreter"
	"github.com/UnnecessaryRain/ironway-core/pkg/network/client"
	"github.com/UnnecessaryRain/ironway-core/pkg/network/protocol"

	log "github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type stdOutSender struct{}

// Send implements the send for this stdout sender
func (s stdOutSender) Send(m protocol.OutgoingMessage) {
	fmt.Println(m)
}

type gameCommand struct {
	player string
}

// Configure sets up the command for gamer
func Configure(app *kingpin.Application) {
	g := &gameCommand{}
	c := app.Command("game", "starts a game").
		Action(g.run)
	// allows us to assume any user for debugging or testing
	c.Flag("assume-user", "username to assume for this game").
		Short('u').
		Required().
		StringVar(&g.player)
}

// run command for game command arg
// runs only the game accepting commands from stdin instead of websockets
func (g *gameCommand) run(c *kingpin.ParseContext) error {
	log.Infoln("Starting standalone game with player", g.player)

	// signal handling to shutdown gracefully
	sigs := make(chan os.Signal)
	stopChan := make(chan struct{})
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		close(stopChan)
		os.Exit(0)
	}()

	metadata := protocol.Metadata{
		Username:  g.player,
		Timestamp: time.Now().Unix(),
	}

	var outWriter stdOutSender

	gameInstance := game.NewGame(outWriter)
	go gameInstance.RunForever(stopChan)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		m := protocol.IncommingMessage(scanner.Text())
		cmd := interpreter.FindCommand(client.Message{
			Metadata: &metadata,
			Message:  &m,
			Sender:   outWriter,
		})
		gameInstance.QueueCommand(outWriter, cmd)
	}

	close(stopChan)

	return nil
}
