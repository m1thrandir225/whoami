import React from 'react'
import { IconLoader } from '@tabler/icons-react'
import { cn } from '@/lib/utils'

interface LoaderProps {
  size?: 'sm' | 'md' | 'lg' | 'xl'
  className?: string
  text?: string
  fullPage?: boolean
}

const sizeMap = {
  sm: 'h-4 w-4',
  md: 'h-6 w-6',
  lg: 'h-8 w-8',
  xl: 'h-12 w-12',
}

const textSizeMap = {
  sm: 'text-sm',
  md: 'text-base',
  lg: 'text-lg',
  xl: 'text-xl',
}

export function Loader({
  size = 'md',
  className,
  text,
  fullPage = false,
}: LoaderProps) {
  const loaderContent = (
    <div
      className={cn(
        'flex flex-col items-center justify-center gap-2',
        className,
      )}
    >
      <IconLoader
        className={cn(sizeMap[size], 'animate-spin text-muted-foreground')}
      />
      {text && (
        <p className={cn(textSizeMap[size], 'text-muted-foreground')}>{text}</p>
      )}
    </div>
  )

  if (fullPage) {
    return (
      <div className="fixed inset-0 flex items-center justify-center bg-background/80 backdrop-blur-sm z-50">
        {loaderContent}
      </div>
    )
  }

  return loaderContent
}

// Specific loader variants for common use cases
export function PageLoader({ text = 'Loading...' }: { text?: string }) {
  return (
    <div className="flex items-center justify-center min-h-[400px]">
      <Loader size="lg" text={text} />
    </div>
  )
}

export function InlineLoader({
  text,
  size = 'sm',
}: {
  text?: string
  size?: 'sm' | 'md'
}) {
  return <Loader size={size} text={text} className="py-4" />
}

export function ButtonLoader({ size = 'sm' }: { size?: 'sm' | 'md' }) {
  return <IconLoader className={cn(sizeMap[size], 'animate-spin')} />
}
