import React, { CSSProperties, FC, Fragment } from "react";
import Highlight, { defaultProps, Language } from "prism-react-renderer";
import theme from "prism-react-renderer/themes/oceanicNext";
import Editor from 'react-simple-code-editor'

interface props {
  value: any
  onValueChange: (code: string) => void
  lang: Language
}

const VirtualEditor: FC<props> = ({ value, onValueChange, lang }) => {

  const styles = {
    root: {
      boxSizing: 'border-box',
      fontFamily: '"Dank Mono", "Fira Code", monospace',
      ...theme.plain
    } as CSSProperties
  }

  const highlight = (code: string) => (
    <Highlight {...defaultProps} theme={theme} code={code} language={lang}>
      {({ className, style, tokens, getLineProps, getTokenProps }) => (
        <Fragment>
          {tokens.map((line, i) => (
            <div {...getLineProps({ line, key: i })}>
              {line.map((token, key) => <span {...getTokenProps({ token, key })} />)}
            </div>
          ))}
        </Fragment>
      )}
    </Highlight>
  )


  return (
    <Editor
      value={value}
      onValueChange={onValueChange}
      highlight={highlight}
      padding={10}
      style={styles.root}
    />
  )
};

export default VirtualEditor;
