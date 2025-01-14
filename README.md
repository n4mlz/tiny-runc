# tiny-runc (WIP)

> **追記**: このプロジェクトはアーカイブされ、開発が終了しました。現在は、同じ目標をかかげて Rust でリライトし、[tiny-youki](https://github.com/n4mlz/tiny-youki) というプロジェクトを開発しています。プロジェクトをアーカイブした理由については [こちら](https://gist.github.com/n4mlz/19533b78a2d83603b42fb30b65936ee0) を参照してください。

低レベルコンテナランタイムである runc の簡易版です (WIP) 。

実用面よりも、コンテナランタイムの動作原理を理解するための参考資料としての活用を目指しています。

以下を目標としています。
- [OCI Runtime Specification](https://github.com/opencontainers/runtime-spec) の必要最低限な部分に準拠
- **rootless に動作**
- シンプルで理解しやすいコード
- コンテナランタイムの自作に関する参考資料としての活用

逆に、以下は目標の範囲としていません。
- OCI Runtime Specification の完全な準拠
- セキュリティの高度な考慮
- プロダクション環境での利用
- 高レベルコンテナランタイムとの連携

## ライセンス

このプロジェクトは [MIT License](LICENSE) の下でライセンスされています。
