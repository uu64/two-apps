import React, { useState } from "react";
import styles from "../styles/Mark.module.css";

const PLUS = "p";
const MINUS = "m";

export type MARK = "p" | "m";

interface Props {
  index: number;
  initValue: MARK;
  onChange: (s: string, i: number) => void;
}

const MarkInput: React.FC<Props> = (props: Props) => {
  const { index, initValue, onChange } = props;
  const [mark, setMark] = useState(initValue);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setMark(e.target.value as MARK);
    onChange(e.target.value, index);
  }

  return (
    <div className={styles.flexcol}>
      <label className={`${styles.mark} ${mark === PLUS ? styles.checked : ""}`}>
        +
        <input
          type="radio"
          name="mark"
          value={PLUS}
          checked={mark === PLUS}
          onChange={handleChange}
        />
      </label>
      <label className={`${styles.mark} ${mark === MINUS ? styles.checked : ""}`}>
        -
        <input
          type="radio"
          name="mark"
          value={MINUS}
          checked={mark === MINUS}
          onChange={handleChange}
        />
      </label>
    </div>
  );
};

export default MarkInput;
