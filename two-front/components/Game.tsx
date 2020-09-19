import React from "react";
import MarkInput, { MARK } from "./MarkInput";
import Number from "./Number";
import styles from "../styles/Game.module.css";

const EQUAL_SYMBOL = "=";
const ANSWER = 2;

interface Props {
  problem: number[];
  answer: MARK[];
  onChange: (answer: MARK, i: number) => void;
}

const Game: React.FC<Props> = (props: Props) => {
  const { problem, answer, onChange } = props;

  const handleChange = (s: MARK, i: number) => {
    onChange(s, i);
  }

  return (
    <>
      <div className={styles.flexrow}><Number number={problem[0]} /></div>
      {answer.map((v, i) => {
        return (
          <div key={i} className={styles.flexrow}>
            <MarkInput
              index={i}
              initValue={v}
              onChange={handleChange}
            />
            <Number number={problem[i + 1]} />
          </div>
        );
      })}
      <div className={styles.flexrow}>
        <div className={styles.equal}>{EQUAL_SYMBOL}</div>
        <Number number={ANSWER} />
      </div>
    </>
  );
};

export default Game;
