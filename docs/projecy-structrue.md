uad-lang/
  go.mod
  README.md
  Makefile

  docs/
    WHITEPAPER.md
    LANGUAGE_SPEC.md
    IR_SPEC.md
    MODEL_LANG_SPEC.md

  cmd/
    uadc/          # compiler: .uad-model / .uad → .uad-IR
      main.go
    uadvm/         # VM runner: 執行 .uad-IR
      main.go
    uadrepl/       # REPL (未來再補)
      main.go

  internal/
    common/
      errors.go
      position.go    # source position (line/column)
      logger.go      # optional
    lexer/
      tokens.go
      lexer.go
    ast/
      core_nodes.go  # .uad-core AST
      model_nodes.go # .uad-model AST
    parser/
      core_parser.go
      model_parser.go
    typer/
      core_types.go
      type_env.go
      type_checker.go
    ir/
      ir.go          # IR type definitions
      builder.go     # AST → IR lowering helpers
      encoder.go     # encode/decode IR (text / binary)
    vm/
      vm.go          # VM core loop
      instr.go       # instruction enum / decoding
      runtime.go     # builtins: math, RNG, domain ops

    model/
      desugar.go     # .uad-model AST -> .uad-core AST
      env.go         # symbol / dataset binding

    runtime/
      stdlib/
        math.uad
        stats.uad
        erh.uad
        security.uad

  examples/
    core/
      hello_world.uad
      erh_basic.uad
    model/
      devsecops_erh.uadmodel
      ransomware_scenario.uadmodel

  scripts/
    dev_setup.sh
    run_examples.sh

  .gitignore
