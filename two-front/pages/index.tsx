import React from "react";
import Head from "next/head";
import Button from "@material-ui/core/Button";
import Snackbar from "@material-ui/core/Snackbar";
import Alert from "@material-ui/lab/Alert";
import { MARK } from "../components/MarkInput";
import Game from "../components/Game";
import styles from "../styles/Home.module.css";

const apiEndpoint = process.env.NEXT_PUBLIC_WS_ENDPOINT;
const level = 5;

interface State {
  message: string;
  openSnackBar: boolean;
  isPlaying: boolean;
  problem: number[];
  answer: MARK[];
}

class Home extends React.Component<{}, State> {
  socket: WebSocket;

  constructor(props: {}) {
    super(props);
    this.state = {
      message: "Please waiting...",
      openSnackBar: true,
      isPlaying: false,
      problem: [],
      answer: [],
    };
  }

  componentDidMount() {
    this.socket = new WebSocket(apiEndpoint);
    this.socket.onopen = () => {
      setTimeout(() => {
        this.startMatching(level);
      }, 3000);
    };
    this.socket.onmessage = this.handleMessage.bind(this);
    this.socket.onclose = this.onDisconnect.bind(this);
  }

  startMatching(level: number) {
    const data = {
      "action": "problem",
      "level": level,
    };
    this.socket.send(JSON.stringify(data));
    setTimeout(() => {
      this.hasNoPlayer();
    }, 60000);
  }

  hasNoPlayer() {
    const { isPlaying } = this.state;
    if (!isPlaying) {
      this.setState({
        message:
          "There is no player. This is disconnected after 3 seconds ...",
      });
      this.disconnect();
    }
  }

  onDisconnect() {
    this.setState({
      message: "This is disconnected.",
    });
    this.openSnackbar();
  }

  handleMessage(ev: MessageEvent): void {
    const data = JSON.parse(ev.data);

    if (!data.message) {
      new Error("Unexpected response");
    }

    switch (data.message) {
      case "PLEASE_WAIT":
        this.waiting();
        break;
      case "START_GAME":
        this.startGame(data.problem);
        break;
      case "WRONG_ANSWER":
        this.isWrongAnswer();
        break;
      case "YOU_WIN":
        this.win();
        break;
      case "YOU_LOSE":
        this.lose();
        break;
      default:
        new Error("Unexpected response");
    }
  }

  waiting() {
    this.setState({
      message: "Looking for a player ...",
    });
    this.openSnackbar();
  }

  startGame(problem: number[]) {
    this.setState({
      message: "Game start !!!",
      isPlaying: true,
      problem: problem,
      answer: Array(problem.length - 1).fill("p"),
    });
    this.openSnackbar();
  }

  isWrongAnswer() {
    this.setState({
      message: "Your answer is wrong :-(",
    });
    this.openSnackbar();
  }

  win() {
    this.setState({
      message: "You win!!! This is disconnected after 3 seconds ...",
    });
    this.openSnackbar();
    this.disconnect();
  }

  lose() {
    this.setState({
      message: "You lose. This is disconnected after 3 seconds ...",
    });
    this.openSnackbar();
    this.disconnect();
  }

  onChange(s: MARK, i: number) {
    const { answer } = this.state;
    answer[i] = s;
    this.setState({
      answer: answer,
    });
  }

  sendAnswer() {
    const { answer } = this.state;
    const data = {
      "action": "solve",
      "answer": answer,
    };
    this.socket.send(JSON.stringify(data));
  }

  disconnect() {
    setTimeout(() => {
      this.socket.close();
    }, 3000);
  }

  openSnackbar() {
    this.setState({
      openSnackBar: true,
    });
  }

  closeSnackbar() {
    this.setState({
      openSnackBar: false,
    });
  }

  render() {
    const { message, openSnackBar, isPlaying, problem, answer } = this.state;
    return (
      <div className={styles.container}>
        <Head>
          <title>Two Apps</title>
          <link rel="icon" href="/favicon.ico" />
        </Head>

        <main className={styles.main}>
          {/* Title */}
          <h1 className={styles.title}>2</h1>
          <p className={styles.description}>Make 2 earlier than another player !</p>

          {/* Game */}
          <div className={styles.grid}>
            {isPlaying
              ? <Game
                  problem={problem}
                  answer={answer}
                  onChange={this.onChange.bind(this)}
                />
              : <div className="loader">Loading...</div>
            }
          </div>
          {isPlaying &&
            <Button
              variant="contained"
              color="primary"
              size="large"
              onClick={this.sendAnswer.bind(this)}
            >
              Answer
            </Button>
          }
        </main>

        {/* message */}
        <Snackbar open={openSnackBar} onClose={this.closeSnackbar.bind(this)}>
          <Alert
            elevation={6}
            variant="filled"
            severity="info"
            onClose={this.closeSnackbar.bind(this)}
          >
            {message}
          </Alert>
        </Snackbar>

        <footer className={styles.footer}>© 2020 uu64</footer>
      </div>
    );
  }
}

export default Home;