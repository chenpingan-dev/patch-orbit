(() => {
  "use strict";
  const $ = (id) => document.getElementById(id);
  const labels = {
    zh: {
      reportLabel: "改动报告",
      filesTab: "变更文件",
      checksTab: "验证清单",
      filesHint: "两个提交之间的文件变更",
      closeDetail: "收起详情",
      offline: "本地分析",
      skip: "跳转到影响地图",
      eyebrow: "GO CHANGE REVIEW",
      headline: "看见改动的连锁反应",
      evidenceKicker: "FOLLOW THE IMPORT",
      edgeDirection: "箭头表示 import 方向",
      manyHops: "4+ 跳依赖",
      moduleScope: "模块复查范围",
      heading: "改动影响地图",
      subtitle: "从一处 Git 改动，沿导入关系找到需要复查的 Go 包",
      import: "打开报告",
      export: "保存报告",
      packages: "包目录",
      affected: "待复查",
      all: "全部包",
      map: "依赖地图",
      reset: "重置选择",
      canvasCaption: "点击包，查看关联路径和导入来源",
      changed: "已改动",
      review: "需复查",
      stable: "未匹配",
      files: "变更文件",
      evidence: "依赖来源",
      checklist: "验证清单",
      download: "下载 ↗",
      copy: "复制检查命令",
      scopeTitle: "分析范围与限制",
      scopeNote:
        "按包导入关系推导，包含测试与构建条件；运行时调用、嵌入资源与外部模块不在首版范围",
      footer: "Patch Orbit · 对比两个已提交的 Git 快照",
      exportJson: "导出 JSON ↗",
      synthetic: "合成数据",
      imported: "导入报告",
      changedPackages: "改动包",
      affectedPackages: "待复查包",
      transitive: "间接影响",
      unaffected: "未匹配包",
      changedDetail: "由 Go 源码变更定位",
      affectedDetail: "沿反向导入关系推导",
      transitiveDetail: "包含两跳及以上依赖",
      unaffectedDetail: "未找到通向改动的导入链",
      search: "搜索包路径",
      empty: "没有匹配的包",
      deleted: "已删除",
      direct: "直接依赖",
      hops: "{n} 跳依赖",
      noChain: "没有找到通向改动包的导入链",
      changedReason: "这个包包含本次修改的 Go 文件",
      globalReason: "依赖清单发生变化，建议检查整个模块",
      globalNote: "go.mod / go.sum 发生变化，所有包纳入复查范围",
      fileCount: "{n} 个 Go 文件",
      focusDescription:
        "建议复查 {n} 个仍存在的包，先运行全量测试，再按地图定位问题",
      notExecuted: "报告生成过程没有执行测试",
      copyFailed: "当前浏览器无法复制，请手动选择命令",
      copied: "已复制",
      saved: "离线报告已保存",
      invalid: "报告格式不正确，请导入 Patch Orbit 生成的 JSON",
      tooLarge: "报告超过 16 MiB 上限",
      shown: "显示 {shown}/{total} 个包",
      notes: "扫描提示",
      new: "新增",
      modified: "修改",
      removed: "删除",
      matching: "{n} 个匹配包",
      downloadTitle: "改动复查清单",
      conservative: "合并改动前后的导入关系，展示保守的复查范围",
      testFiles: "测试文件",
      noPackages: "快照中没有可解析的 Go 包",
      globalCommand: "先运行",
      focusedCommand: "缩小定位范围",
      basedOn: "导入依据",
      inspect: "查看",
      directImports: "直接导入对照",
      sourceFilter: "筛选导入来源",
      allSources: "全部来源",
      baseOnly: "仅改动前",
      headOnly: "仅改动后",
      bothSnapshots: "两边都有",
      before: "改动前",
      after: "改动后",
      sourceCount: "{shown}/{total} 个文件来源",
      noSources: "没有符合筛选条件的来源",
      noImports: "未解析到对本地包的直接导入",
      sourceScope: "来源筛选不改变地图的复查范围",
      sourceHelp: "如何对照",
      sourceHelpText:
        "按同一文件和 import 路径配对，行号移动仍算两边都有；仅前或仅后表示只在该快照解析到，扫描跳过的文件也可能影响结果",
      fullPath: "查看完整关联路径",
    },
    en: {
      reportLabel: "Change report",
      filesTab: "Files",
      checksTab: "Checks",
      filesHint: "Files changed between the two commits",
      closeDetail: "Close details",
      offline: "Local analysis",
      skip: "Skip to impact map",
      eyebrow: "GO CHANGE REVIEW",
      headline: "See where a change travels",
      evidenceKicker: "FOLLOW THE IMPORT",
      edgeDirection: "Arrows follow imports",
      manyHops: "4+ import hops",
      moduleScope: "Module review scope",
      heading: "Change impact map",
      subtitle: "Follow a Git change through Go packages, one import at a time",
      import: "Open a report",
      export: "Save report",
      packages: "Packages",
      affected: "Review scope",
      all: "All packages",
      map: "Import map",
      reset: "Reset selection",
      canvasCaption: "Select a package to inspect its import path",
      changed: "Changed",
      review: "Review",
      stable: "Unmatched",
      files: "Changed files",
      evidence: "Import evidence",
      checklist: "Verification checklist",
      download: "Download ↗",
      copy: "Copy test commands",
      scopeTitle: "Analysis scope & limitations",
      scopeNote:
        "Package imports, including tests and build-tag variants; runtime calls, embedded assets and external modules are outside v0.1",
      footer: "Patch Orbit · Compare two committed Git snapshots",
      exportJson: "Export JSON ↗",
      synthetic: "Synthetic data",
      imported: "Imported report",
      changedPackages: "Changed packages",
      affectedPackages: "Packages to review",
      transitive: "Transitive dependents",
      unaffected: "Unmatched packages",
      changedDetail: "Mapped from changed Go sources",
      affectedDetail: "Follow reverse import edges",
      transitiveDetail: "Two or more import hops",
      unaffectedDetail: "No import chain to a changed package",
      search: "Search package paths",
      empty: "No matching packages",
      deleted: "Deleted",
      direct: "Direct import",
      hops: "{n} import hops",
      noChain: "No import chain to a changed package was found",
      changedReason: "This package contains changed Go source files",
      globalReason: "The dependency manifest changed; review the entire module",
      globalNote: "go.mod / go.sum changed; every package is in scope",
      fileCount: "Go files: {n}",
      focusDescription:
        "Review {n} surviving packages. Run all tests first, then use the map to narrow the search",
      notExecuted: "Generating this report did not execute tests",
      copyFailed: "Clipboard is unavailable; select and copy the commands",
      copied: "Copied",
      saved: "Offline report saved",
      invalid: "Invalid report; import JSON generated by Patch Orbit",
      tooLarge: "Report exceeds the 16 MiB limit",
      shown: "Showing {shown}/{total} packages",
      notes: "Scan notes",
      new: "Added",
      modified: "Modified",
      removed: "Deleted",
      matching: "{n} matching packages",
      downloadTitle: "Change review checklist",
      conservative:
        "Union of base/head imports gives a conservative review scope",
      testFiles: "Test import",
      noPackages: "No parseable Go packages in the snapshots",
      globalCommand: "Start here",
      focusedCommand: "Narrow the search",
      basedOn: "Import evidence",
      inspect: "Inspect",
      directImports: "Direct import comparison",
      sourceFilter: "Filter import sources",
      allSources: "All sources",
      baseOnly: "Base only",
      headOnly: "Head only",
      bothSnapshots: "Both snapshots",
      before: "Base",
      after: "Head",
      sourceCount: "{shown}/{total} file sources",
      noSources: "No sources match this filter",
      noImports: "No direct imports of local packages were parsed",
      sourceScope: "Source filters keep the map's review scope unchanged",
      sourceHelp: "How sources are paired",
      sourceHelpText:
        "Pair the same file and import path, even when line numbers move. Base/head only means parsed in that snapshot only; skipped files can also affect the result",
      fullPath: "Inspect the full impact path",
    },
  };
  let locale = "zh";
  try {
    locale =
      localStorage.getItem("independent-tools.locale") === "en" ? "en" : "zh";
  } catch {}
  const t = (key, values = {}) =>
    Object.entries(values).reduce(
      (text, [name, value]) => text.replaceAll(`{${name}}`, String(value)),
      labels[locale][key] ?? key,
    );
  const esc = (value) =>
    String(value ?? "").replace(
      /[&<>"']/g,
      (c) =>
        ({
          "&": "&amp;",
          "<": "&lt;",
          ">": "&gt;",
          '"': "&quot;",
          "'": "&#39;",
        })[c],
    );
  const num = (value) =>
    Number(value).toLocaleString(locale === "zh" ? "zh-CN" : "en-US");
  function validate(data) {
    if (
      data?.schema !== "patch-orbit/v1" ||
      !Array.isArray(data.nodes) ||
      data.nodes.length > 6000 ||
      !Array.isArray(data.edges) ||
      !Array.isArray(data.changes) ||
      !data.summary ||
      !Array.isArray(data.notes)
    )
      throw Error(t("invalid"));
    const ids = new Set();
    for (const n of data.nodes) {
      if (
        !n ||
        typeof n.id !== "string" ||
        ids.has(n.id) ||
        !Array.isArray(n.files) ||
        !Array.isArray(n.chain)
      )
        throw Error(t("invalid"));
      ids.add(n.id);
    }
    const integer = (v) => Number.isSafeInteger(v) && v >= 0;
    if (
      !["title", "module", "base", "head", "test_command"].every(
        (k) => typeof data[k] === "string",
      ) ||
      !data.notes.every((n) => typeof n === "string") ||
      ![
        "packages",
        "changed_files",
        "changed_packages",
        "affected_packages",
        "direct",
        "transitive",
        "unaffected",
      ].every((k) => integer(data.summary[k]))
    )
      throw Error(t("invalid"));
    for (const n of data.nodes) {
      if (
        !n.files.every((p) => typeof p === "string") ||
        !n.chain.every((p) => ids.has(p)) ||
        typeof n.changed !== "boolean" ||
        typeof n.present !== "boolean" ||
        !(n.distance === null || integer(n.distance))
      )
        throw Error(t("invalid"));
    }
    for (const edge of data.edges) {
      if (
        !edge ||
        !ids.has(edge.from) ||
        !ids.has(edge.to) ||
        !Array.isArray(edge.evidence) ||
        !edge.evidence.every(
          (e) =>
            e &&
            ["file", "import", "snapshot"].every(
              (k) => typeof e[k] === "string",
            ) &&
            ["base", "head"].includes(e.snapshot) &&
            integer(e.line) &&
            e.line > 0 &&
            (e.test === undefined || typeof e.test === "boolean"),
        )
      )
        throw Error(t("invalid"));
    }
    if (
      !data.changes.every(
        (c) => c && typeof c.path === "string" && typeof c.status === "string",
      )
    )
      throw Error(t("invalid"));
    return data;
  }
  const initial = (data) =>
    [...data.nodes]
      .filter((n) => n.distance !== null)
      .sort((a, b) => b.distance - a.distance)[0]?.id ?? data.nodes[0]?.id;
  let data = validate(JSON.parse($("report-data").textContent)),
    selected = initial(data),
    scope = "affected",
    query = "",
    activeView = "map",
    sourceFilter = "all",
    inspectorOpen = false;
  let toastTimer;
  function toast(text) {
    $("toast").textContent = text;
    $("toast").classList.add("visible");
    clearTimeout(toastTimer);
    toastTimer = setTimeout(() => $("toast").classList.remove("visible"), 2600);
  }
  function status(n) {
    return !n.present
      ? "deleted"
      : n.changed
        ? "changed"
        : n.distance !== null
          ? "review"
          : "stable";
  }
  function metricCards() {
    const s = data.summary;
    $("metrics").innerHTML = [
      ["changedPackages", s.changed_packages, "changedDetail"],
      ["affectedPackages", s.affected_packages, "affectedDetail"],
      ["transitive", s.transitive, "transitiveDetail"],
      ["unaffected", s.unaffected, "unaffectedDetail"],
    ]
      .map(
        ([label, value, detail]) =>
          `<div class="metric" title="${esc(t(detail))}"><strong>${esc(num(value))}</strong><span>${esc(t(label))}</span></div>`,
      )
      .join("");
  }
  function renderView() {
    ["map", "files", "checks"].forEach((view) => {
      $(view + "-view").hidden = activeView !== view;
      $(view + "-tab").classList.toggle("active", activeView === view);
      $(view + "-tab").setAttribute(
        "aria-pressed",
        String(activeView === view),
      );
    });
    const visible = inspectorOpen && activeView === "map";
    $("inspector").hidden = !visible;
    $("map-workspace").classList.toggle("inspector-open", visible);
    $("reset").hidden = activeView !== "map";
  }
  function openDetail() {
    inspectorOpen = true;
    activeView = "map";
    renderView();
    renderSelection();
  }
  function closeDetail() {
    inspectorOpen = false;
    renderView();
    const selectedNode = [
      ...$("graph").querySelectorAll("[data-node-index]"),
    ].find((node) => node.getAttribute("aria-label") === selected);
    (selectedNode ?? $("map-tab")).focus({ preventScroll: true });
  }
  function candidates() {
    return data.nodes.filter(
      (n) =>
        (scope === "all" || n.distance !== null) &&
        n.id.toLowerCase().includes(query),
    );
  }
  function packageList() {
    const nodes = candidates();
    $("package-count").textContent = `${nodes.length}/${data.nodes.length}`;
    $("package-list").innerHTML =
      nodes
        .map(
          (n) =>
            `<button class="package-row ${n.id === selected ? "active" : ""}" data-package="${esc(n.id)}" aria-pressed="${n.id === selected}"><span class="package-symbol ${n.changed ? "changed" : ""}"><svg viewBox="0 0 20 20" aria-hidden="true"><path d="m3 6 7-4 7 4v8l-7 4-7-4V6Zm0 0 7 4 7-4m-7 4v8"/></svg></span><span><strong>${esc(n.id)}</strong><small>${esc(!n.present ? t("deleted") : n.changed ? t("changed") : n.distance === null ? t("stable") : n.distance === 1 ? t("direct") : t("hops", { n: n.distance }))}</small></span></button>`,
        )
        .join("") || `<div class="empty">${esc(t("empty"))}</div>`;
    $("package-list")
      .querySelectorAll("[data-package]")
      .forEach(
        (button) =>
          (button.onclick = () => {
            selected = button.dataset.package;
            openDetail();
          }),
      );
    $("package-note").textContent = t("matching", { n: nodes.length });
    $("affected").setAttribute("aria-pressed", String(scope === "affected"));
    $("all").setAttribute("aria-pressed", String(scope === "all"));
    $("affected").classList.toggle("active", scope === "affected");
    $("all").classList.toggle("active", scope === "all");
  }
  function draw() {
    const current = data.nodes.find((n) => n.id === selected);
    const chain = new Set(current?.chain ?? []);
    let list = candidates();
    if (!query && current && !list.some((n) => n.id === current.id))
      list = [current, ...list];
    const visible = [...list]
      .sort(
        (a, b) =>
          Number(chain.has(b.id)) - Number(chain.has(a.id)) ||
          (a.distance ?? 999) - (b.distance ?? 999) ||
          a.id.localeCompare(b.id),
      )
      .slice(0, 80);
    const groups = new Map();
    for (const node of visible) {
      const level = node.distance === null ? 99 : Math.min(4, node.distance);
      if (!groups.has(level)) groups.set(level, []);
      groups.get(level).push(node);
    }
    const levels = [...groups.keys()].sort((a, b) => b - a);
    const width = Math.max(700, levels.length * 178 + 32);
    const rows = Math.max(
      1,
      ...[...groups.values()].map((nodes) => nodes.length),
    );
    const height = Math.max(360, rows * 86 + 130);
    const positions = new Map();
    let svg = `<defs><marker id="arrow" markerUnits="userSpaceOnUse" markerWidth="6" markerHeight="6" refX="5" refY="3" orient="auto"><path d="M0,0 L6,3 L0,6" fill="#b7c5db"/></marker><marker id="hot-arrow" markerUnits="userSpaceOnUse" markerWidth="6" markerHeight="6" refX="5" refY="3" orient="auto"><path d="M0,0 L6,3 L0,6" fill="#5484e2"/></marker></defs>`;
    levels.forEach((level, column) => {
      const x = (width - (levels.length - 1) * 178) / 2 + column * 178;
      const nodes = groups.get(level);
      const label =
        level === 99
          ? t("stable")
          : level === 0
            ? t(data.module_changed ? "moduleScope" : "changed")
            : level === 1
              ? t("direct")
              : level === 4
                ? t("manyHops")
                : t("hops", { n: level });
      svg += `<text class="column-label" x="${x}" y="35" text-anchor="middle">${esc(label)}</text><path class="column-line" d="M${x} 54 V${height - 23}"/>`;
      nodes.forEach((node, i) =>
        positions.set(node.id, {
          x,
          y: 80 + ((i + 0.5) * (height - 115)) / nodes.length,
        }),
      );
    });
    const chainLinks = new Set(
      (current?.chain ?? [])
        .slice(0, -1)
        .map((id, i) => id + "\0" + current.chain[i + 1]),
    );
    for (const edge of data.edges) {
      const a = positions.get(edge.from),
        b = positions.get(edge.to);
      if (!a || !b) continue;
      const hot = chainLinks.has(edge.from + "\0" + edge.to);
      const sameColumn = a.x === b.x;
      const middle = sameColumn ? a.x + 113 : (a.x + b.x) / 2;
      const destination = sameColumn ? b.x + 74 : b.x - 76;
      svg += `<path class="graph-edge ${hot ? "highlight" : ""}" d="M${a.x + 74},${a.y} C${middle},${a.y} ${middle},${b.y} ${destination},${b.y}" marker-end="url(#${hot ? "hot-arrow" : "arrow"})"/>`;
    }
    visible.forEach((node, i) => {
      const p = positions.get(node.id),
        parts = node.id.split("/");
      const title = parts.pop() || node.id;
      const parent = parts.join("/") || ".";
      const shorten = (text, length) =>
        text.length > length ? text.slice(0, length - 1) + "…" : text;
      const state = status(node);
      svg += `<g class="graph-node ${state} ${node.id === selected ? "selected" : ""}" data-node-index="${i}" transform="translate(${p.x},${p.y})" tabindex="0" role="button" aria-pressed="${node.id === selected}" aria-label="${esc(node.id)}"><title>${esc(node.id)} · ${esc(t(state))}</title><rect class="node-card" x="-73" y="-31" width="146" height="62" rx="9"/><circle class="node-dot" cx="-57" cy="-10" r="3"/><text class="node-name" x="-47" y="-6">${esc(shorten(title, 14))}</text><text class="node-path" x="-57" y="12">${esc(shorten(parent, 24))}</text></g>`;
    });
    $("graph").setAttribute("viewBox", `0 0 ${width} ${height}`);
    $("graph").style.minWidth = `${Math.max(660, width * 0.86)}px`;
    $("graph").innerHTML = svg;
    $("graph")
      .querySelectorAll("[data-node-index]")
      .forEach((node) => {
        const choose = () => {
          const hadFocus = document.activeElement === node;
          selected = visible[Number(node.dataset.nodeIndex)].id;
          openDetail();
          const replacement = [
            ...$("graph").querySelectorAll("[data-node-index]"),
          ].find((element) => element.getAttribute("aria-label") === selected);
          if (hadFocus && replacement)
            replacement.focus({ preventScroll: true });
        };
        node.onclick = choose;
        node.onkeydown = (event) => {
          if (event.key === "Enter" || event.key === " ") {
            event.preventDefault();
            choose();
          }
        };
      });
    $("visible-count").textContent = t("shown", {
      shown: visible.length,
      total: list.length,
    });
    $("graph-empty").hidden = visible.length > 0;
    $("graph-empty").textContent = t(
      data.nodes.length ? "empty" : "noPackages",
    );
  }
  function snapshotState(base, head) {
    if (base.length && head.length) {
      return "both";
    }
    return base.length ? "base" : "head";
  }
  const sourceLabel = (state) =>
    t({ base: "baseOnly", head: "headOnly", both: "bothSnapshots" }[state]);
  function groupEvidence(evidence) {
    const groups = new Map();
    for (const proof of evidence) {
      // A moved import keeps its identity; line numbers belong to each snapshot.
      const key = JSON.stringify([
        proof.file,
        proof.import,
        Boolean(proof.test),
      ]);
      if (!groups.has(key)) {
        groups.set(key, {
          file: proof.file,
          import: proof.import,
          test: Boolean(proof.test),
          base: new Set(),
          head: new Set(),
        });
      }
      groups.get(key)[proof.snapshot].add(proof.line);
    }
    return [...groups.values()]
      .map((group) => {
        const base = [...group.base].sort((a, b) => a - b);
        const head = [...group.head].sort((a, b) => a - b);
        return { ...group, base, head, state: snapshotState(base, head) };
      })
      .sort(
        (a, b) =>
          a.file.localeCompare(b.file) || a.import.localeCompare(b.import),
      );
  }
  function sourceHTML(group) {
    const positions = [
      ["base", "before", data.base],
      ["head", "after", data.head],
    ]
      .map(
        ([snapshot, label, commit]) =>
          `<div class="source-position ${group[snapshot].length ? "" : "absent"}"><span title="${esc(commit)}">${esc(t(label))}</span><strong>${group[snapshot].length ? group[snapshot].map((line) => `L${num(line)}`).join(", ") : "—"}</strong></div>`,
      )
      .join("");
    return `<div class="import-evidence" data-source-state="${group.state}"><div class="evidence-source"><strong>${esc(group.file)}</strong><span class="snapshot-tag ${group.state}">${esc(sourceLabel(group.state))}</span></div><div class="source-positions">${positions}</div><code>import "${esc(group.import)}"</code>${group.test ? `<span class="test-source">${esc(t("testFiles"))}</span>` : ""}</div>`;
  }
  function directImports(n) {
    const imports = data.edges
      .filter((edge) => edge.from === n.id)
      .map((edge) => ({ ...edge, groups: groupEvidence(edge.evidence) }))
      .sort((a, b) => a.to.localeCompare(b.to));
    if (!imports.length) {
      return `<section class="direct-imports"><h4>${esc(t("directImports"))}</h4><p class="source-empty">${esc(t("noImports"))}</p></section>`;
    }

    const total = imports.reduce(
      (count, edge) => count + edge.groups.length,
      0,
    );
    const visible = imports
      .map((edge) => ({
        ...edge,
        visible: edge.groups.filter(
          (group) => sourceFilter === "all" || group.state === sourceFilter,
        ),
      }))
      .filter((edge) => edge.visible.length);
    const shown = visible.reduce(
      (count, edge) => count + edge.visible.length,
      0,
    );
    const options = [
      ["all", "allSources"],
      ["base", "baseOnly"],
      ["head", "headOnly"],
      ["both", "bothSnapshots"],
    ]
      .map(
        ([value, label]) =>
          `<option value="${value}" ${sourceFilter === value ? "selected" : ""}>${esc(t(label))}</option>`,
      )
      .join("");
    const cards = visible
      .map((edge, index) => {
        const state = snapshotState(
          edge.evidence.filter((proof) => proof.snapshot === "base"),
          edge.evidence.filter((proof) => proof.snapshot === "head"),
        );
        return `<details class="dependency-card" data-import-target="${esc(edge.to)}" ${index === 0 ? "open" : ""}><summary><strong>${esc(edge.to)}</strong><span class="snapshot-tag ${state}">${esc(sourceLabel(state))}</span></summary>${edge.visible.map(sourceHTML).join("")}</details>`;
      })
      .join("");
    return `<section class="direct-imports"><h4>${esc(t("directImports"))}</h4><label class="source-filter" for="source-filter"><span>${esc(t("sourceFilter"))}</span><select id="source-filter">${options}</select></label><p class="source-count" role="status">${esc(t("sourceCount", { shown, total }))}</p>${cards || `<p class="source-empty">${esc(t("noSources"))}</p>`}<p class="source-scope">${esc(t("sourceScope"))}</p><details class="source-help"><summary>${esc(t("sourceHelp"))}</summary><p>${esc(t("sourceHelpText"))}</p></details></section>`;
  }
  function impactPath(n) {
    let html = `<details class="path-evidence"><summary>${esc(t("fullPath"))}</summary><div class="chain-list">`;
    n.chain.forEach((id, index) => {
      html += `<div class="chain-step">${esc(id)}</div>`;
      if (index < n.chain.length - 1) {
        const edge = data.edges.find(
          (candidate) =>
            candidate.from === id && candidate.to === n.chain[index + 1],
        );
        html += `<div class="chain-arrow">↓ ${esc(t("basedOn"))}</div>`;
        html += groupEvidence(edge?.evidence ?? [])
          .map(sourceHTML)
          .join("");
      }
    });
    return html + "</div></details>";
  }
  function detail() {
    const n = data.nodes.find((n) => n.id === selected);
    if (!n) {
      $("node-status").textContent = "";
      $("node-detail").innerHTML =
        `<div class="empty">${esc(t("empty"))}</div>`;
      return;
    }
    const state = status(n);
    $("node-status").className = `pill ${state}`;
    $("node-status").textContent = t(state);
    let html = `<div class="node-overview"><h3>${esc(n.id)}</h3><p>${esc(t("fileCount", { n: n.files.length }))}</p><div class="node-files">${n.files.slice(0, 5).map(esc).join("<br>")}</div></div>`;
    if (n.changed)
      html += `<div class="node-empty">${esc(t("changedReason"))}</div>`;
    else if (data.module_changed)
      html += `<div class="node-empty">${esc(t("globalReason"))}</div>`;
    else if (n.chain.length < 2)
      html += `<div class="node-empty">${esc(t("noChain"))}</div>`;
    html += directImports(n);
    if (!n.changed && !data.module_changed && n.chain.length >= 2) {
      html += impactPath(n);
    }
    $("node-detail").innerHTML = html;
    const filter = $("source-filter");
    if (filter) {
      filter.onchange = () => {
        sourceFilter = filter.value;
        detail();
        $("source-filter").focus({ preventScroll: true });
      };
    }
  }
  function commands() {
    const targets = data.nodes
      .filter((n) => n.present && n.distance !== null)
      .map((n) => n.id)
      .filter((id) => /^[a-zA-Z0-9._/-]+$/.test(id));
    const focused =
      targets.length && targets.length < 30
        ? "go test " +
          targets.map((id) => (id === "." ? "." : "./" + id)).join(" ")
        : "go test ./...";
    return {
      focused,
      full: "go test ./...",
      count: data.nodes.filter((n) => n.present && n.distance !== null).length,
    };
  }
  function plan() {
    const c = commands();
    $("test-plan").innerHTML =
      `<p class="test-description">${esc(t("focusDescription", { n: c.count }))}</p><p class="test-description"><strong>${esc(t("globalCommand"))}</strong></p><div class="test-command">${esc(c.full)}</div>${c.focused !== c.full ? `<p class="test-description" style="margin-top:10px"><strong>${esc(t("focusedCommand"))}</strong></p><div class="test-command">${esc(c.focused)}</div>` : ""}<p class="test-description" style="margin-top:10px">${esc(t("notExecuted"))}</p>`;
  }
  function changes() {
    $("change-count").textContent =
      `${Math.min(30, data.changes.length)}/${data.changes.length}`;
    $("change-list").innerHTML =
      data.changes
        .slice(0, 30)
        .map(
          (c) =>
            `<div class="change-row"><span>${esc(c.path)}</span><b class="change-status ${["A", "D", "M"].includes(c.status) ? c.status : ""}">${esc(c.status)} · ${esc(t(c.status === "A" ? "new" : c.status === "D" ? "removed" : "modified"))}</b></div>`,
        )
        .join("") || `<div class="empty">${esc(t("empty"))}</div>`;
  }
  function renderSelection() {
    packageList();
    draw();
    detail();
  }
  function render() {
    document.documentElement.lang = locale === "zh" ? "zh-CN" : "en";
    document.title = `Patch Orbit · ${t("heading")}`;
    document
      .querySelectorAll("[data-i18n]")
      .forEach((node) => (node.textContent = t(node.dataset.i18n)));
    $("language").textContent = locale === "zh" ? "English" : "中文";
    $("language").setAttribute(
      "aria-label",
      locale === "zh" ? "Switch to English" : "切换为中文",
    );
    $("search").placeholder = t("search");
    $("search").setAttribute("aria-label", t("search"));
    $("sample-badge").textContent = t(
      data.synthetic ? "synthetic" : "imported",
    );
    $("report-title").textContent = data.title_i18n?.[locale] ?? data.title;
    $("module").textContent = data.module;
    $("base").textContent = data.base.slice(0, 7);
    $("head").textContent = data.head.slice(0, 7);
    $("close-detail").setAttribute("aria-label", t("closeDetail"));
    $("graph").setAttribute("aria-label", t("map"));
    $("inspector").setAttribute("aria-label", t("evidence"));
    metricCards();
    renderView();
    renderSelection();
    changes();
    plan();
    const notes = data.notes.filter(
      (n) =>
        !n.startsWith("package imports only") &&
        !n.startsWith("go.mod/go.sum changed"),
    );
    $("notes").innerHTML =
      (data.module_changed ? `<p>${esc(t("globalNote"))}</p>` : "") +
      `<p>${esc(t("conservative"))}</p>` +
      (notes.length
        ? `<details><summary>${esc(t("notes"))} (${notes.length})</summary>${notes.map((n) => `<p>${esc(n)}</p>`).join("")}</details>`
        : "");
  }
  function download(content, type, filename) {
    const url = URL.createObjectURL(new Blob([content], { type }));
    const a = document.createElement("a");
    a.href = url;
    a.download = filename;
    a.click();
    setTimeout(() => URL.revokeObjectURL(url), 1000);
  }
  $("language").onclick = () => {
    locale = locale === "zh" ? "en" : "zh";
    try {
      localStorage.setItem("independent-tools.locale", locale);
    } catch {}
    render();
  };
  $("search").oninput = (e) => {
    query = e.target.value.toLowerCase();
    renderSelection();
  };
  $("affected").onclick = () => {
    scope = "affected";
    renderSelection();
  };
  $("all").onclick = () => {
    scope = "all";
    renderSelection();
  };
  document.querySelectorAll("[data-view]").forEach((button) => {
    button.onclick = () => {
      activeView = button.dataset.view;
      renderView();
    };
  });
  $("close-detail").onclick = closeDetail;
  document.addEventListener("keydown", (event) => {
    if (event.key === "Escape" && inspectorOpen) {
      closeDetail();
    }
  });
  $("reset").onclick = () => {
    query = "";
    $("search").value = "";
    selected = data.nodes.find((n) => n.changed)?.id ?? data.nodes[0]?.id;
    scope = "affected";
    sourceFilter = "all";
    inspectorOpen = false;
    renderView();
    renderSelection();
  };
  $("copy").onclick = async () => {
    const c = commands();
    try {
      await navigator.clipboard.writeText(
        c.full + (c.focused !== c.full ? "\n" + c.focused : ""),
      );
      toast(t("copied"));
    } catch {
      toast(t("copyFailed"));
    }
  };
  $("markdown").onclick = () => {
    const c = commands();
    let text =
      `# ${t("downloadTitle")}\n\n${data.module}\n\n${data.base} → ${data.head}\n\n## ${t("files")}\n\n` +
      data.changes.map((x) => `- ${x.status}: ${x.path}`).join("\n");
    text +=
      `\n\n## ${t("affected")}\n\n` +
      data.nodes
        .filter((n) => n.distance !== null)
        .map((n) => `- ${n.id}: ${n.chain.join(" → ")}`)
        .join("\n");
    text += `\n\n## ${t("checklist")}\n\n${c.full}\n${c.focused !== c.full ? c.focused : ""}\n\n${t("scopeNote")}\n`;
    download(text, "text/markdown", "patch-orbit-checklist.md");
  };
  $("json-export").onclick = () =>
    download(
      JSON.stringify(data, null, 2),
      "application/json",
      "patch-orbit-report.json",
    );
  $("export").onclick = () => {
    const clone = document.documentElement.cloneNode(true);
    clone.querySelector("#report-data").textContent = JSON.stringify(
      data,
    ).replaceAll("<", "\\u003c");
    clone.querySelector("#toast").classList.remove("visible");
    download(
      "<!doctype html>\n" + clone.outerHTML,
      "text/html",
      "patch-orbit-report.html",
    );
    toast(t("saved"));
  };
  $("import").onclick = () => $("file").click();
  $("file").onchange = async () => {
    const file = $("file").files[0];
    if (!file) return;
    try {
      if (file.size > 16 * 1024 * 1024) throw Error(t("tooLarge"));
      data = validate(JSON.parse(await file.text()));
      selected = data.nodes.find((n) => n.changed)?.id ?? data.nodes[0]?.id;
      query = "";
      scope = "affected";
      sourceFilter = "all";
      $("search").value = "";
      inspectorOpen = false;
      activeView = "map";
      render();
    } catch (error) {
      toast(error.message || t("invalid"));
    } finally {
      $("file").value = "";
    }
  };
  render();
})();
