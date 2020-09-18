import React, { useState } from "react";
import Mark from "./Mark";
import Number from "./Number";
import styles from "../styles/Game.module.css";

interface Props {
  problem: number[];
  onChange: (answer: string[]) => void;
}

const Game: React.FC<Props> = (props: Props) => {
  const { problem, onChange } = props;
  const [answer, setAnswer] = useState(Array(problem.length - 1).fill("p"));

  const handleChange = (s: string, i: number) => {
    answer[i] = s;
    setAnswer(answer);
    onChange(answer);
  }

  return (
    <>
      <div className={styles.flexrow}><Number number={problem[0]} /></div>
      {problem.slice(1, problem.length).map((v, i) => {
        return (
          <div key={i} className={styles.flexrow}>
            <Mark
              index={i}
              onChange={handleChange}
            />
            <Number number={v} />
          </div>
        );
      })}
      <div className={styles.flexrow}>
        <div className={styles.equal}>=</div>
        <Number number={2} />
      </div>
    </>
  );
};

export default Game;
