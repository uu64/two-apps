import React from "react";
import Head from "next/head";
import Game from "../components/Game";
import styles from "../styles/Home.module.css";

const apiEndpoint = "wss://nubmwv3y2g.execute-api.ap-northeast-1.amazonaws.com/dev";
const level = 4;

interface State {
  message: string;
  isPlaying: boolean;
  problem: number[];
}

class Home extends React.Component<{}, State> {
  socket: WebSocket;

  constructor(props: {}) {
    super(props);
    this.state = {
      message: "Please waiting...",
      isPlaying: false,
      problem: [],
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
  }

  startMatching(level: number) {
    const data = {
      "action": "problem",
      "level": level,
    };
    this.socket.send(JSON.stringify(data));
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
        // this.dispatchEvent(new Event("win"));
        break;
      case "YOU_LOSE":
        // this.dispatchEvent(new Event("lose"));
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

  startGame(problem: any) {
    this.setState({
      message: "Game start!!!",
      isPlaying: true,
      problem: problem,
    });
  }

  isWrongAnswer() {
    this.setState({
      message: "Your answer is wrong :-(",
    });
  }

  onChange(answer: string[]) {
    console.log(answer);
  }

  render() {
    const { message, isPlaying, problem } = this.state;
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
              ? <Game problem={problem} onChange={this.onChange} />
              : <div className="loader">Loading...</div>
            }
          </div>
        </main>

        <footer className={styles.footer}>© 2020 uu64</footer>
      </div>
    );
  }
}

export default Home;