import React from "react";
import styles from "../styles/Number.module.css";

interface Props {
  number: number;
}

const Number: React.FC<Props> = (props: Props) => {
  const { number } = props;
  return (
    <div>
      <span className={styles.number}>
        {number}
      </span>
    </div>
  );
};

export default Number;
