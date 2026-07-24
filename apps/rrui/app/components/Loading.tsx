'use client';

export default function Loading() {

  return (
    <div className="h-screen flex items-center justify-center">
      <svg className="h-6 w-6 animate-spin" viewBox="0 0 100 100">
        <circle
          fill="none"
          stroke-width="10"
          className="stroke-current opacity-40"
          cx="50"
          cy="50"
          r="40"
        />
        <circle
          fill="none"
          stroke-width="10"
          className="stroke-current"
          stroke-dasharray="250"
          stroke-dashoffset="210"
          cx="50"
          cy="50"
          r="40"
        />
      </svg>
    </div>

  )
}

