import React from "react";
import Head from "next/head";
import { MARK } from "../components/MarkInput";
import Game from "../components/Game";
import styles from "../styles/Home.module.css";

const apiEndpoint = "wss://nubmwv3y2g.execute-api.ap-northeast-1.amazonaws.com/dev";
const level = 5;

interface State {
  message: string;
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
    this.socket.onclose = this.onClose.bind(this);
  }

  startMatching(level: number) {
    const data = {
      "action": "problem",
      "level": level,
    };
    this.socket.send(JSON.stringify(data));
    setTimeout(() => {
      this.hasNoChallenger();
    }, 60000);
  }

  answer() {
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
      message: "Please waiting...",
    });
  }

  startGame(problem: number[]) {
    this.setState({
      message: "Game start!!!",
      isPlaying: true,
      problem: problem,
      answer: Array(problem.length - 1).fill("p"),
    });
  }

  isWrongAnswer() {
    this.setState({
      message: "Your answer is wrong :-(",
    });
  }

  win() {
    this.setState({
      message: "You win!!! This is disconnected after 3 seconds...",
    });
    this.disconnect();
  }

  lose() {
    this.setState({
      message: "You lose. This is disconnected after 3 seconds...",
    });
    this.disconnect();
  }

  hasNoChallenger() {
    const { isPlaying } = this.state;
    if (!isPlaying) {
      this.setState({
        message: "There is no challenger. This is disconnected after 3 seconds...",
      });
      this.disconnect();
    }
  }

  onChange(s: MARK, i: number) {
    const { answer } = this.state;
    answer[i] = s;
    this.setState({
      answer: answer,
    });
  }

  onClose() {
    this.setState({
      message: "This is disconnected.",
    });
  }

  render() {
    const { message, isPlaying, problem, answer } = this.state;
    return (
      <div className={styles.container}>
        <Head>
          <title>Two Apps</title>
          <link rel="icon" href="/favicon.ico" />
        </Head>

        <main className={styles.main}>
          {/* Title */}
          <h1 className={styles.title}>2</h1>
          <p className={styles.description}>Lets make 2 !</p>

          {/* Message */}
          <div className={styles.log}>
            <code className={styles.code}>{message}</code>
          </div>

          {/* Game */}
          <div className={styles.grid}>
            {isPlaying
              ? <Game problem={problem} answer={answer} onChange={this.onChange.bind(this)} />
              : <div className="loader">Loading...</div>
            }
          </div>
          {isPlaying &&
            <div className={styles.button} onClick={this.answer.bind(this)}>
              Answer
            </div>
          }
        </main>

        <footer className={styles.footer}>© 2020 uu64</footer>
      </div>
    );
  }
}

export default Home;