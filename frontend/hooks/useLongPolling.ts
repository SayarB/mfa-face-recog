import { useCallback, useEffect, useRef, useState } from "react";

export const useLongPolling = (url: string, options: any, enabled: boolean) => {
  const [data, setData] = useState<{
    isComplete: boolean;
    isSuccess: boolean;
    isFailure: boolean;
  } | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<any>(null);
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const fetchData = useCallback(async () => {
    try {
      setIsLoading(true);
      const response = await fetch(url, options);
      const json = await response.json();
      setData(json);
    } catch (error) {
      setError(error);
    } finally {
      setIsLoading(false);
    }
  }, [url, options]);

  let i = 0;
  useEffect(() => {
    if (i === 1) return;
    if (!enabled) return;
    if (!!timeoutRef.current) return;
    timeoutRef.current = setInterval(fetchData, 1000);
    i++;
    return () => {
      if (timeoutRef.current) {
        clearInterval(timeoutRef.current);
      }
    };
  }, [enabled]);

  return { data, isLoading, error };
};
