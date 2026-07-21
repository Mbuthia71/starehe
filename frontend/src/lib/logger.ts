// Logging utilities for the Siohioma SACCO system

export enum LogLevel {
  DEBUG = 0,
  INFO = 1,
  WARN = 2,
  ERROR = 3,
  FATAL = 4,
}

export interface LogEntry {
  timestamp: string;
  level: LogLevel;
  message: string;
  context?: Record<string, any>;
  userId?: string;
  sessionId?: string;
}

class Logger {
  private logs: LogEntry[] = [];
  private maxLogs = 1000;
  private minLevel: LogLevel = LogLevel.INFO;

  constructor() {
    // Send logs to server in production
    if (typeof window === 'undefined') {
      setInterval(() => this.flushLogs(), 30000); // Flush every 30s
    }
  }

  setMinLevel(level: LogLevel): void {
    this.minLevel = level;
  }

  private log(level: LogLevel, message: string, context?: Record<string, any>): void {
    if (level < this.minLevel) return;

    const entry: LogEntry = {
      timestamp: new Date().toISOString(),
      level,
      message,
      context,
      userId: this.getUserId(),
      sessionId: this.getSessionId(),
    };

    this.logs.push(entry);

    // Keep only last maxLogs entries
    if (this.logs.length > this.maxLogs) {
      this.logs = this.logs.slice(-this.maxLogs);
    }

    // Console output in development
    if (import.meta.env.DEV) {
      const levelName = LogLevel[level];
      const emoji = this.getLevelEmoji(level);
      console.log(`[${emoji} ${levelName}] ${message}`, context || '');
    }
  }

  private getLevelEmoji(level: LogLevel): string {
    switch (level) {
      case LogLevel.DEBUG: return '🔍';
      case LogLevel.INFO: return 'ℹ️';
      case LogLevel.WARN: return '⚠️';
      case LogLevel.ERROR: return '❌';
      case LogLevel.FATAL: return '💀';
      default: return '📝';
    }
  }

  private getUserId(): string | undefined {
    if (typeof window === 'undefined') return undefined;
    try {
      return localStorage.getItem('siohioma_clientId') || undefined;
    } catch {
      return undefined;
    }
  }

  private getSessionId(): string | undefined {
    if (typeof window === 'undefined') return undefined;
    try {
      return sessionStorage.getItem('session_id') || undefined;
    } catch {
      return undefined;
    }
  }

  debug(message: string, context?: Record<string, any>): void {
    this.log(LogLevel.DEBUG, message, context);
  }

  info(message: string, context?: Record<string, any>): void {
    this.log(LogLevel.INFO, message, context);
  }

  warn(message: string, context?: Record<string, any>): void {
    this.log(LogLevel.WARN, message, context);
  }

  error(message: string, context?: Record<string, any>): void {
    this.log(LogLevel.ERROR, message, context);
  }

  fatal(message: string, context?: Record<string, any>): void {
    this.log(LogLevel.FATAL, message, context);
  }

  getLogs(level?: LogLevel): LogEntry[] {
    if (level !== undefined) {
      return this.logs.filter(log => log.level >= level);
    }
    return [...this.logs];
  }

  clearLogs(): void {
    this.logs = [];
  }

  async flushLogs(): Promise<void> {
    if (this.logs.length === 0) return;

    try {
      // In production, send logs to server
      if (!import.meta.env.DEV && typeof window !== 'undefined') {
        await fetch('/api/logs', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ logs: this.logs }),
        });
      }
      this.logs = [];
    } catch (error) {
      console.error('Failed to flush logs:', error);
    }
  }
}

export const logger = new Logger();

// Error boundary logging
export function logError(error: Error, context?: Record<string, any>): void {
  logger.error(error.message, {
    stack: error.stack,
    name: error.name,
    ...context,
  });
}

// API error logging
export function logAPIError(endpoint: string, error: any, context?: Record<string, any>): void {
  logger.error(`API Error: ${endpoint}`, {
    status: error.status,
    statusText: error.statusText,
    ...context,
  });
}
