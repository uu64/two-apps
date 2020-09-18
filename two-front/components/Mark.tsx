import React, { useState } from "react";
import styles from "../styles/Mark.module.css";

interface Props {
  index: number;
  onChange: (s: string, i: number) => void;
}

const Mark: React.FC<Props> = (props: Props) => {
  const { index, onChange } = props;
  const [mark, setMark] = useState("p");

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setMark(e.target.value);
    onChange(e.target.value, index);
  }

  return (
    <div className={styles.flexcol}>
      <label className={`${styles.mark} ${mark === "p" ? styles.checked : ""}`}>
        +
        <input
          type="radio"
          name="mark"
          value="p"
          checked={mark === "p"}
          onChange={handleChange}
        />
      </label>
      <label className={`${styles.mark} ${mark === "m" ? styles.checked : ""}`}>
        -

        <input
          type="radio"
          name="mark"
          value="m"
          checked={mark === "m"}
          onChange={handleChange}
        />
      </label>
    </div>
  );
};

export default Mark;
