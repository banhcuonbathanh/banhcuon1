"use client";

import { useState, useEffect } from "react";
import { ToastContainer, toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";

interface TimerProps {
  initialMinutes: number;
  initialSeconds: number;
  onTimerEnd?: () => void;
}

const Timer: React.FC<TimerProps> = ({
  initialMinutes,
  initialSeconds,
  onTimerEnd
}) => {
  const [timerMinutes, setTimerMinutes] = useState(initialMinutes);
  const [timerSeconds, setTimerSeconds] = useState(initialSeconds);
  const [isTimerRunning, setIsTimerRunning] = useState(false);
  const [isPaused, setIsPaused] = useState(false);

  useEffect(() => {
    let interval: any;

    if (isTimerRunning) {
      interval = setInterval(() => {
        setTimerSeconds((prevSeconds) => {
          if (prevSeconds === 0) {
            setTimerMinutes((prevMinutes) => {
              if (prevMinutes === 0 && prevSeconds === 0) {
                toast.info("Time's up!", {
                  position: "top-right",
                  autoClose: 5000,
                  hideProgressBar: false,
                  closeOnClick: true,
                  pauseOnHover: true,
                  draggable: true,
                  progress: undefined,
                  theme: "light"
                });
                setIsTimerRunning(false);
                setIsPaused(false);
                onTimerEnd?.();
                return 0;
              }
              return prevMinutes - 1;
            });
            return 59;
          } else {
            return prevSeconds - 1;
          }
        });
      }, 1000);
    }

    return () => clearInterval(interval);
  }, [isTimerRunning, isPaused, onTimerEnd]);

  const startTimer = () => {
    setIsTimerRunning(true);
    setIsPaused(false);
  };

  const stopTimer = () => {
    setIsTimerRunning(false);
    setIsPaused(false);
    setTimerMinutes(initialMinutes);
    setTimerSeconds(initialSeconds);
  };

  const pauseTimer = () => {
    setIsPaused(!isPaused);
  };

  const updateTimer = () => {
    setTimerMinutes(initialMinutes);
    setTimerSeconds(initialSeconds);
  };

  return (
    <div>
      <ToastContainer />
      <div className="text-center text-2xl font-bold mb-8">
        Time Remaining: {timerMinutes.toString().padStart(2, "0")}:
        {timerSeconds.toString().padStart(2, "0")}
      </div>
      <div className="mb-8">
        <button
          className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded ml-4"
          onClick={startTimer}
        >
          Start Timer
        </button>
        <button
          className="bg-yellow-500 hover:bg-yellow-700 text-white font-bold py-2 px-4 rounded ml-4"
          onClick={pauseTimer}
        >
          {isPaused ? "Resume" : "Pause"}
        </button>
        <button
          className="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded ml-4"
          onClick={stopTimer}
        >
          Stop Timer
        </button>
        <button
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded ml-4"
          onClick={updateTimer}
        >
          Update Timer
        </button>
      </div>
    </div>
  );
};

export default Timer;
