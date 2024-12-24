import React from 'react';

interface GridContainerProps {
  children: React.ReactNode;
  className?: string;
}

const GridContainer = ({ children, className = '' }: GridContainerProps) => {
  return (
    <div className="container mx-auto px-4 py-8">
      <div className={`grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 ${className}`}>
        {children}
      </div>
    </div>
  );
};

export default GridContainer;