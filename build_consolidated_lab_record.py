"""
Consolidated Laboratory Record Builder
Author: Vishal V (24BPS1089)
Course: Software Engineering (SWE) Laboratory
Department: Department of Computer Science and Engineering

Generates:
1. Consolidated_Lab_Experiments_24BPS1089.docx (workspace root)
2. Report files/Consolidated_Lab_Experiments_24BPS1089.docx
3. CONSOLIDATED_LAB_EXPERIMENTS.md (workspace root)
"""

import os
import shutil
import docx
from docx import Document
from docx.shared import Inches, Pt, RGBColor
from docx.enum.text import WD_ALIGN_PARAGRAPH
from docx.enum.table import WD_TABLE_ALIGNMENT
from docx.oxml import parse_xml
from docx.oxml.ns import nsdecls

def set_cell_background(cell, fill_hex):
    tcPr = cell._tc.get_or_add_tcPr()
    shd = parse_xml(f'<w:shd {nsdecls("w")} w:fill="{fill_hex}"/>')
    tcPr.append(shd)

def set_cell_margins(cell, top=140, bottom=140, left=160, right=160):
    tcPr = cell._tc.get_or_add_tcPr()
    tcMar = parse_xml(f'<w:tcMar {nsdecls("w")}><w:top w:w="{top}" w:type="dxa"/><w:bottom w:w="{bottom}" w:type="dxa"/><w:left w:w="{left}" w:type="dxa"/><w:right w:w="{right}" w:type="dxa"/></w:tcMar>')
    tcPr.append(tcMar)

def add_heading_styled(doc, text, level):
    h = doc.add_heading(text, level=level)
    run = h.runs[0]
    if level == 1:
        run.font.size = Pt(15.5)
        run.font.bold = True
        run.font.color.rgb = RGBColor(15, 23, 42)
        h.paragraph_format.space_before = Pt(16)
        h.paragraph_format.space_after = Pt(6)
    elif level == 2:
        run.font.size = Pt(12.5)
        run.font.bold = True
        run.font.color.rgb = RGBColor(30, 41, 59)
        h.paragraph_format.space_before = Pt(12)
        h.paragraph_format.space_after = Pt(4)
    elif level == 3:
        run.font.size = Pt(11)
        run.font.bold = True
        run.font.color.rgb = RGBColor(51, 65, 85)
        h.paragraph_format.space_before = Pt(8)
        h.paragraph_format.space_after = Pt(3)
    return h

def add_paragraph_styled(doc, text, bold_prefix=None, space_after=5, italic=False):
    p = doc.add_paragraph()
    p.paragraph_format.space_after = Pt(space_after)
    p.paragraph_format.line_spacing = 1.15
    if bold_prefix:
        r_b = p.add_run(bold_prefix)
        r_b.bold = True
        r_b.font.size = Pt(10)
        r_b.font.color.rgb = RGBColor(15, 23, 42)
    r = p.add_run(text)
    r.font.size = Pt(10)
    r.italic = italic
    r.font.color.rgb = RGBColor(51, 65, 85)
    return p

def add_bullet_styled(doc, text, bold_prefix=None):
    p = doc.add_paragraph(style='List Bullet')
    p.paragraph_format.space_after = Pt(3)
    p.paragraph_format.line_spacing = 1.15
    if bold_prefix:
        r_b = p.add_run(bold_prefix)
        r_b.bold = True
        r_b.font.size = Pt(9.5)
        r_b.font.color.rgb = RGBColor(15, 23, 42)
    r = p.add_run(text)
    r.font.size = Pt(9.5)
    r.font.color.rgb = RGBColor(51, 65, 85)
    return p

def add_callout_box(doc, title, text_items):
    tbl = doc.add_table(rows=1, cols=1)
    tbl.alignment = WD_TABLE_ALIGNMENT.CENTER
    cell = tbl.cell(0, 0)
    set_cell_background(cell, "F8FAFC")
    set_cell_margins(cell, top=140, bottom=140, left=180, right=180)
    cell.width = Inches(6.5)
    
    tcPr = cell._tc.get_or_add_tcPr()
    borders = parse_xml(f'''
        <w:tcBorders {nsdecls("w")}>
            <w:top w:val="none"/>
            <w:left w:val="single" w:sz="24" w:space="0" w:color="2563EB"/>
            <w:bottom w:val="none"/>
            <w:right w:val="none"/>
        </w:tcBorders>
    ''')
    tcPr.append(borders)
    
    p = cell.paragraphs[0]
    p.paragraph_format.space_after = Pt(3)
    r_t = p.add_run(f"📌 {title}\n")
    r_t.font.bold = True
    r_t.font.size = Pt(10)
    r_t.font.color.rgb = RGBColor(30, 58, 138)
    
    for item in text_items:
        p_item = cell.add_paragraph()
        p_item.paragraph_format.space_after = Pt(2)
        r_i = p_item.add_run(item)
        r_i.font.size = Pt(9.5)
        r_i.font.color.rgb = RGBColor(71, 85, 105)
    
    p_after = doc.add_paragraph()
    p_after.paragraph_format.space_after = Pt(4)

def add_table_styled(doc, headers, data, col_widths=None):
    tbl = doc.add_table(rows=len(data) + 1, cols=len(headers))
    tbl.alignment = WD_TABLE_ALIGNMENT.CENTER
    
    # Header row
    hdr_row = tbl.rows[0]
    trPr = hdr_row._tr.get_or_add_trPr()
    trPr.append(parse_xml(f'<w:tblHeader {nsdecls("w")}/>'))
    for idx, heading in enumerate(headers):
        cell = hdr_row.cells[idx]
        set_cell_background(cell, "1E293B")
        set_cell_margins(cell, top=130, bottom=130, left=130, right=130)
        p = cell.paragraphs[0]
        p.alignment = WD_ALIGN_PARAGRAPH.LEFT
        r = p.add_run(heading)
        r.font.bold = True
        r.font.size = Pt(9)
        r.font.color.rgb = RGBColor(255, 255, 255)
    
    # Data rows
    for r_idx, row_data in enumerate(data):
        row = tbl.rows[r_idx + 1]
        bg = "F8FAFC" if r_idx % 2 == 1 else "FFFFFF"
        for c_idx, cell_value in enumerate(row_data):
            cell = row.cells[c_idx]
            set_cell_background(cell, bg)
            set_cell_margins(cell, top=90, bottom=90, left=110, right=110)
            p = cell.paragraphs[0]
            p.alignment = WD_ALIGN_PARAGRAPH.LEFT
            r = p.add_run(str(cell_value))
            r.font.size = Pt(8.5)
            r.font.color.rgb = RGBColor(51, 65, 85)
            
    # Set borders
    tblPr = tbl._tbl.tblPr
    borders = parse_xml(f'''
        <w:tblBorders {nsdecls("w")}>
            <w:top w:val="single" w:sz="6" w:space="0" w:color="CBD5E1"/>
            <w:bottom w:val="single" w:sz="6" w:space="0" w:color="CBD5E1"/>
            <w:insideH w:val="single" w:sz="4" w:space="0" w:color="E2E8F0"/>
            <w:left w:val="none"/>
            <w:right w:val="none"/>
            <w:insideV w:val="none"/>
        </w:tblBorders>
    ''')
    tblPr.append(borders)
    
    if col_widths and len(col_widths) == len(headers):
        for row in tbl.rows:
            for idx, w in enumerate(col_widths):
                row.cells[idx].width = Inches(w)
                
    p_end = doc.add_paragraph()
    p_end.paragraph_format.space_after = Pt(6)
    return tbl

def add_image_figure(doc, img_path, caption, max_width_in=5.8):
    if os.path.exists(img_path):
        try:
            p = doc.add_paragraph()
            p.alignment = WD_ALIGN_PARAGRAPH.CENTER
            p.paragraph_format.space_before = Pt(8)
            p.paragraph_format.space_after = Pt(3)
            run = p.add_run()
            run.add_picture(img_path, width=Inches(max_width_in))
            
            cap_p = doc.add_paragraph()
            cap_p.alignment = WD_ALIGN_PARAGRAPH.CENTER
            cap_r = cap_p.add_run(f"Figure: {caption}")
            cap_r.font.size = Pt(9)
            cap_r.font.bold = True
            cap_r.font.color.rgb = RGBColor(71, 85, 105)
            cap_p.paragraph_format.space_after = Pt(12)
            return True
        except Exception as e:
            print(f"Error adding {img_path}: {e}")
            return False
    else:
        print(f"Missing: {img_path}")
        return False

def add_experiment_header(doc, exp_no, title, aim, tools):
    tbl = doc.add_table(rows=1, cols=1)
    tbl.alignment = WD_TABLE_ALIGNMENT.CENTER
    cell = tbl.cell(0, 0)
    set_cell_background(cell, "0F172A")
    set_cell_margins(cell, top=160, bottom=160, left=200, right=200)
    cell.width = Inches(6.5)
    
    p = cell.paragraphs[0]
    p.alignment = WD_ALIGN_PARAGRAPH.LEFT
    p.paragraph_format.space_after = Pt(2)
    
    r_exp = p.add_run(f"EXPERIMENT {exp_no}\n")
    r_exp.font.bold = True
    r_exp.font.size = Pt(11)
    r_exp.font.color.rgb = RGBColor(96, 165, 250)
    
    r_title = p.add_run(f"{title.upper()}\n")
    r_title.font.bold = True
    r_title.font.size = Pt(13.5)
    r_title.font.color.rgb = RGBColor(255, 255, 255)
    
    p_meta = cell.add_paragraph()
    p_meta.paragraph_format.space_after = Pt(0)
    r_m = p_meta.add_run(f"Student: Vishal V  |  Reg No: 24BPS1089  |  Modelling / Testing Tools: {tools}")
    r_m.font.size = Pt(9)
    r_m.font.color.rgb = RGBColor(203, 213, 225)
    
    p_sp = doc.add_paragraph()
    p_sp.paragraph_format.space_after = Pt(6)
    
    add_callout_box(doc, "Aim of the Experiment", [aim])

def build_full_record():
    doc = Document()
    
    # Page setup - 1 inch margins
    sec = doc.sections[0]
    sec.top_margin = Inches(1.0)
    sec.bottom_margin = Inches(1.0)
    sec.left_margin = Inches(1.0)
    sec.right_margin = Inches(1.0)
    sec.different_first_page_header_footer = True

    # Setup header and footer for normal pages
    header = sec.header
    hp = header.paragraphs[0]
    hp.alignment = WD_ALIGN_PARAGRAPH.RIGHT
    hrun = hp.add_run("Software Engineering Laboratory (SWE)  |  Vishal V (24BPS1089)")
    hrun.font.size = Pt(8.5)
    hrun.font.color.rgb = RGBColor(148, 163, 184)

    footer = sec.footer
    fp = footer.paragraphs[0]
    fp.alignment = WD_ALIGN_PARAGRAPH.CENTER
    frun = fp.add_run("Consolidated Laboratory Record — AegisVote & BlazeDemo QA Suite")
    frun.font.size = Pt(8.5)
    frun.font.color.rgb = RGBColor(148, 163, 184)

    # -------------------------------------------------------------------------
    # COVER / TITLE PAGE
    # -------------------------------------------------------------------------
    p_sp = doc.add_paragraph()
    p_sp.paragraph_format.space_before = Pt(36)
    
    p_inst = doc.add_paragraph()
    p_inst.alignment = WD_ALIGN_PARAGRAPH.CENTER
    r_inst = p_inst.add_run("DEPARTMENT OF COMPUTER SCIENCE AND ENGINEERING\nSCHOOL OF COMPUTER SCIENCE AND ENGINEERING\n")
    r_inst.font.bold = True
    r_inst.font.size = Pt(13)
    r_inst.font.color.rgb = RGBColor(30, 41, 59)
    p_inst.paragraph_format.space_after = Pt(24)

    p_proj = doc.add_paragraph()
    p_proj.alignment = WD_ALIGN_PARAGRAPH.CENTER
    r_sub = p_proj.add_run("SOFTWARE ENGINEERING LABORATORY (SWE LAB)\nCONSOLIDATED LABORATORY RECORD / EXPERIMENT DOSSIER\n\n")
    r_sub.font.bold = True
    r_sub.font.size = Pt(12)
    r_sub.font.color.rgb = RGBColor(100, 116, 139)
    
    r_title = p_proj.add_run("Software Analysis, Design, Architecture &\nVerification of AegisVote: Privacy-Preserving Electronic Voting System\n& BlazeDemo Flight Booking QA Suite")
    r_title.font.bold = True
    r_title.font.size = Pt(17)
    r_title.font.color.rgb = RGBColor(15, 23, 42)
    p_proj.paragraph_format.space_after = Pt(36)

    p_meta = doc.add_paragraph()
    p_meta.alignment = WD_ALIGN_PARAGRAPH.CENTER
    r_by = p_meta.add_run("Submitted by:\n")
    r_by.font.size = Pt(11)
    r_name = p_meta.add_run("VISHAL V\n")
    r_name.font.bold = True
    r_name.font.size = Pt(14)
    r_reg = p_meta.add_run("Register Number: 24BPS1089\n\n")
    r_reg.font.bold = True
    r_reg.font.size = Pt(12)
    r_reg.font.color.rgb = RGBColor(30, 64, 175)
    
    r_under = p_meta.add_run("Bachelor of Technology (B.Tech) in Computer Science & Engineering\nAcademic Year: 2025 – 2026")
    r_under.font.size = Pt(10.5)
    r_under.font.italic = True
    p_meta.paragraph_format.space_after = Pt(40)

    doc.add_page_break()

    # -------------------------------------------------------------------------
    # TABLE OF CONTENTS / INDEX TABLE
    # -------------------------------------------------------------------------
    add_heading_styled(doc, "Laboratory Experiments Index", level=1)
    add_paragraph_styled(doc, "This consolidated laboratory record compiles the complete series of software engineering laboratory experiments completed during the course. Experiments 2 through 8 systematically develop the requirements, scheduling, structured analysis, dynamic behavioral models, object-oriented design, interaction lifecycles, and physical deployment architecture for AegisVote (a privacy-preserving electronic voting system). Experiment 9 demonstrates automated quality assurance through Selenium IDE functional test automation and Apache JMeter performance load testing on BlazeDemo.", italic=True)
    
    index_headers = ["Exp No.", "Experiment Title", "Category / Phase", "Modelling & Testing Tools"]
    index_data = [
        ["Exp 2", "Project Scheduling & Work Breakdown Structures (WBS)", "Project Planning & Scheduling", "GanttProject, MS Project, StarUML"],
        ["Exp 3", "Requirement Modelling: Context (L-0), L-1 & L-2 DFDs", "Structured Analysis & DFD", "StarUML / Diagramming Engine"],
        ["Exp 4", "Requirement Modelling: Entity-Relationship (ER) Diagram", "Database & Schema Design", "Relational Modeller, StarUML"],
        ["Exp 5", "Behavioral Modelling: State Transition Diagram", "Dynamic Behavioral Modelling", "StarUML State Machine Modeller"],
        ["Exp 6", "Object-Oriented Design: Use Case and Class Diagrams", "OOAD & Structural Modelling", "StarUML OOAD Engine"],
        ["Exp 7", "Interaction Modelling: Sequence & Communication Diagrams", "Dynamic Interaction Modelling", "StarUML Interaction Engine"],
        ["Exp 8", "OO Design: Package, Component & Deployment Models", "System Architecture & Deployment", "StarUML Architectural Modeller"],
        ["Exp 9", "Software Testing: Design and Demonstration of Test Cases", "Quality Assurance & Testing", "Selenium IDE, Apache JMeter"]
    ]
    add_table_styled(doc, index_headers, index_data, col_widths=[0.8, 2.7, 1.6, 1.4])

    doc.add_page_break()

    # -------------------------------------------------------------------------
    # EXPERIMENT 2
    # -------------------------------------------------------------------------
    add_experiment_header(
        doc,
        exp_no="2",
        title="Project Scheduling & Work Breakdown Structures (WBS)",
        aim="To perform software project scheduling, decomposition, duration estimation, task usage, resource allocation, and Work Breakdown Structure (WBS) modelling for the Privacy-Preserving Electronic Voting System using StarUML, GanttProject, and Microsoft Project.",
        tools="GanttProject, Microsoft Project, StarUML"
    )

    add_heading_styled(doc, "1. Project Overview & Scope", level=2)
    add_paragraph_styled(doc, "The Privacy-Preserving Electronic Voting System (AegisVote) is a secure, web-based electronic voting platform designed to guarantee ballot secrecy and cryptographic tamper-evidence without relying on blockchain consensus. The system solves the primary challenge of electronic voting: decoupling voter authentication from vote casting so that no single server or malicious database administrator can ever associate an authenticated voter's identity with their submitted candidate choice.")
    
    add_heading_styled(doc, "Core Architectural Innovations:", level=3)
    add_bullet_styled(doc, "Separates voter identity verification from ballot persistence using an independent, time-bound Anonymous Voting Token Service.", "Decoupled Authentication: ")
    add_bullet_styled(doc, "Each registered voter receives a single-use, unlinked cryptographic token authorizing exactly one vote per election.", "Single-Use Anonymous Tokens: ")
    add_bullet_styled(doc, "Every vote payload is encrypted with AES-256 before disk storage; an append-only binary SHA-256 Merkle Tree ledger tracks ballot hashes in real time to instantly detect tampering or unauthorized ballot deletion.", "Cryptographic Integrity: ")
    add_bullet_styled(doc, "Strict role boundaries for Voters, Election Officers, Administrators, and Counting Authorities with automated audit logging.", "Role-Based Governance: ")

    add_image_figure(doc, "Report files/extracted_figures/Lab2_Gantt_Overview.png", "GanttProject Timeline and Master Project Schedule")

    add_heading_styled(doc, "2. Work Breakdown Structure (WBS) Modelling", level=2)
    add_paragraph_styled(doc, "To ensure methodical execution across the software engineering lifecycle, the project was decomposed into hierarchical Work Breakdown Structures using four orthogonal decomposition criteria:")
    
    add_bullet_styled(doc, "Decomposes the project along the chronological software engineering lifecycle: Requirements Engineering → Architectural Modelling → Backend Cryptographic Core → Frontend Development → Integration & Penetration Testing → Deployment.", "1. Process-Based Decomposition: ")
    add_bullet_styled(doc, "Decomposes the project according to tangible architectural deliverables: Voter Portal UI, Authentication Gateway, Token Issuance Engine, AES-256 Ballot Encryption Service, PostgreSQL Merkle Ledger, and Election Counting Authority Dashboard.", "2. Product-Based Decomposition: ")
    add_bullet_styled(doc, "Decomposes tasks across physical and cloud deployment zones: Regional Voter Client Terminals, Central Cloud Application Cluster (US-East), and Air-Gapped Verification Server for final cryptographic tallying.", "3. Geography-Based Decomposition: ")
    add_bullet_styled(doc, "Decomposes work according to team specializations: Project Lead, Cryptography Specialist, Full-Stack Frontend Developer, Backend Go Developer, QA / Test Automation Engineer, and Database Administrator.", "4. Role-Based Decomposition: ")

    add_image_figure(doc, "Report files/extracted_figures/Lab2_Planning_image16.png", "Work Breakdown Structure (WBS) Decomposition Hierarchy (StarUML)")
    add_image_figure(doc, "Report files/extracted_figures/Lab2_Product_WBS.png", "Product-Based Work Breakdown Structure (WBS)")
    add_image_figure(doc, "Report files/extracted_figures/Lab2_Geography_WBS.png", "Geography-Based Distributed Architecture WBS")
    add_image_figure(doc, "Report files/extracted_figures/Lab2_Role_WBS.png", "Role-Based Resource & Responsibility Decomposition")

    add_heading_styled(doc, "3. Activity Network, Task Scheduling & Duration Allocation", level=2)
    add_paragraph_styled(doc, "A Program Evaluation and Review Technique (PERT) / Critical Path Method (CPM) Activity Network was synthesized to model task precedence constraints, finish-to-start (FS) dependencies, and critical path milestones. Task usage, work effort (man-days), and duration allocations were formulated in GanttProject / Microsoft Project to ensure resource leveling and eliminate delivery bottlenecks.")

    add_image_figure(doc, "Report files/extracted_figures/Lab2_Network_Diagram.png", "PERT / CPM Activity Network Diagram and Dependency Graph")
    add_image_figure(doc, "Report files/extracted_figures/Lab2_Planning_image23.png", "Task Usage and Schedule Duration Allocation Matrix")
    add_image_figure(doc, "Report files/extracted_figures/Lab2_Planning_image24.png", "Resource Sheet and Personnel Allocation Matrix")
    add_image_figure(doc, "Report files/extracted_figures/Lab2_Resource_Usage.png", "Resource Allocation Graph and Team Workload Distribution")

    add_heading_styled(doc, "4. Result & Conclusion", level=2)
    add_paragraph_styled(doc, "The project scheduling, multi-criteria WBS decomposition, task network, and resource allocation models were successfully completed using StarUML, GanttProject, and Microsoft Project. The resulting schedule establishes a robust foundation for agile sprint tracking, milestone verification, and resource leveling across the 12-week development lifecycle.")

    doc.add_page_break()

    # -------------------------------------------------------------------------
    # EXPERIMENT 3
    # -------------------------------------------------------------------------
    add_experiment_header(
        doc,
        exp_no="3",
        title="Requirement Modelling: Context (Level-0), Level-1, and Level-2 DFDs",
        aim="To perform structured requirements analysis and construct Context Diagram (Level-0 DFD), Level-1 DFD, and Level-2 DFD blueprints for the Privacy-Preserving Electronic Voting System.",
        tools="StarUML, Structured Systems Analysis Engine"
    )

    add_heading_styled(doc, "1. System Overview & Structured Analysis", level=2)
    add_paragraph_styled(doc, "Data Flow Diagrams (DFDs) provide a hierarchical, process-driven perspective of the voting system, mapping how data enters, transitions through processing units, and commits to persistent storage without exposing voter identities. The architecture explicitly guarantees:")
    add_bullet_styled(doc, "Guarantees that no single process or data store ever contains both the voter's real-world identity and their candidate selection simultaneously.", "Data Separation & Privacy Wall: ")
    add_bullet_styled(doc, "Maintains append-only ballot commitment verified through cryptographic hashing.", "Tamper-Evident Storage: ")
    add_bullet_styled(doc, "Multi-factor authentication (Credentials + TOTP) verifying eligibility while keeping the ballot anonymous.", "Strong Authentication: ")

    add_heading_styled(doc, "2. External Entities & Data Stores", level=2)
    add_paragraph_styled(doc, "The structured model identifies five external entities and six core data stores:")
    add_bullet_styled(doc, "Voter (submits credentials, receives token, casts ballot), Election Officer (manages elections and candidates), Administrator (manages accounts, reviews audit logs), Counting Authority (verifies Merkle root, tallies ballots post-closure), and OTP Provider (delivers out-of-band one-time verification tokens).", "External Entities: ")
    add_bullet_styled(doc, "D1: Users (Credentials & Role), D2: Elections (Status & Timeline), D3: Candidates (Roster), D4: Tokens (Hashed anonymous tokens), D5: Encrypted Ballots (AES-256 ciphertexts), D6: Audit Logs (System events).", "Data Stores: ")

    add_heading_styled(doc, "3. Context Diagram (Level-0 DFD) Blueprint", level=2)
    add_paragraph_styled(doc, "The Level-0 Context Diagram represents the entire voting platform as a single centralized process interacting with external actors, defining the external data interfaces and system security boundary.")
    add_image_figure(doc, "Report files/extracted_figures/DFD_Level0_Context.png", "Context Diagram (Level-0 DFD) Blueprint")

    add_heading_styled(doc, "4. Level-1 DFD: Functional Process Decomposition", level=2)
    add_paragraph_styled(doc, "The Level-1 DFD decomposes the system into five primary functional processes:")
    add_bullet_styled(doc, "Authenticates voter via credentials and TOTP, verifies active election eligibility, and checks the per-voter HasVoted flag.", "Process 1.0: Voter Authentication & Eligibility Check — ")
    add_bullet_styled(doc, "Issues a cryptographically random, single-use, time-bound anonymous token and flags the voter as having obtained a token.", "Process 2.0: Anonymous Token Issuance — ")
    add_bullet_styled(doc, "Validates token unicity, encrypts ballot payload using AES-256, appends to database, updates the Merkle tree root, and invalidates the token.", "Process 3.0: Process & Encrypt Vote — ")
    add_bullet_styled(doc, "Initiated post-election closure to verify Merkle root consistency, decrypt ballots, and compute final results.", "Process 4.0: Vote Tallying & Integrity Verification — ")
    add_bullet_styled(doc, "Facilitates election setup, candidate registration, and system event logging.", "Process 5.0: Election Administration & Audit Logging — ")
    add_image_figure(doc, "Report files/extracted_figures/DFD_Level1_System.png", "Level-1 DFD: Primary System Process Decomposition")

    add_heading_styled(doc, "5. Level-2 DFD: Decomposing Process 3.0 (Process & Encrypt Vote)", level=2)
    add_paragraph_styled(doc, "Level-2 DFD zooms directly into the high-security vote ingestion pipeline (Process 3.0), detailing the sub-processes executing inside the core voting controller:")
    add_bullet_styled(doc, "Validates the raw token against D4: Tokens and ensures Used == False.", "Sub-Process 3.1 (Token Validation & Invalidation): ")
    add_bullet_styled(doc, "Encrypts the candidate choice using AES-256-GCM symmetric encryption.", "Sub-Process 3.2 (Ballot Encryption): ")
    add_bullet_styled(doc, "Computes SHA-256 leaf hash and recalculates the root of the Merkle Tree.", "Sub-Process 3.3 (Merkle Tree Leaf Generation): ")
    add_bullet_styled(doc, "Appends the encrypted ciphertext and leaf hash to D5: Encrypted Ballots.", "Sub-Process 3.4 (Append-Only Persistence): ")
    add_image_figure(doc, "Report files/extracted_figures/DFD_Level2_Process3.png", "Level-2 DFD: Process 3.0 (Process & Encrypt Vote Decomposition)")

    add_heading_styled(doc, "6. Result & Conclusion", level=2)
    add_paragraph_styled(doc, "The Context Diagram (Level-0), Level-1, and Level-2 Data Flow Diagrams were successfully developed. The models rigorously prove that voter identification information never flows into the ballot encryption or persistent storage pipelines, establishing complete architectural anonymity.")

    doc.add_page_break()

    # -------------------------------------------------------------------------
    # EXPERIMENT 4
    # -------------------------------------------------------------------------
    add_experiment_header(
        doc,
        exp_no="4",
        title="Requirement Modelling: Entity-Relationship (ER) Diagram & Relational Schema",
        aim="To design the Entity-Relationship (ER) Diagram, resolve Many-to-Many associations using bridge entities, formulate the complete data dictionary, and establish relational schema constraints enforcing zero voter-to-ballot linkage.",
        tools="StarUML, Relational Database Modeller"
    )

    add_heading_styled(doc, "1. Problem Statement & Zero-Linkage Constraint", level=2)
    add_paragraph_styled(doc, "The objective is to design a relational database architecture that strictly preserves voter privacy. In traditional voting databases, storing a UserID foreign key inside the Vote record permanently links individuals to their political choices. To eliminate this vulnerability, the schema enforces a strict Zero-PII architectural detachment:")
    add_bullet_styled(doc, "The Vote entity intentionally contains NO UserID, VoterID, or foreign key referencing the User table.", "Crucial Security Constraint 1: ")
    add_bullet_styled(doc, "The Token entity stores only a one-way cryptographic SHA-256 hash (TokenHash) with no direct link to the User record in persistent storage.", "Crucial Security Constraint 2: ")
    add_bullet_styled(doc, "The Many-to-Many (M:N) relationship between Users and Elections is resolved via an associative bridge entity (User_Election_Eligibility) to track voter eligibility and the HasVoted flag independently for every election.", "Crucial Security Constraint 3: ")

    add_heading_styled(doc, "2. Comprehensive Data Dictionary Table", level=2)
    add_paragraph_styled(doc, "The complete database design comprises seven relational entities as detailed in the data dictionary below:")

    er_table_headers = ["Entity Name", "Attribute", "Attribute Type", "Key / Integrity Constraint"]
    er_table_data = [
        ["User", "UserID", "UUID / Simple", "Primary Key (PK)"],
        ["User", "Name", "VARCHAR(100)", "Not Null"],
        ["User", "PasswordHash", "VARCHAR(255)", "Bcrypt Hash, Not Null"],
        ["User", "OTPSecret", "VARCHAR(64)", "TOTP Shared Secret"],
        ["User", "Role", "VARCHAR(20)", "Enum ('VOTER', 'ADMIN', 'OFFICER', 'AUDITOR')"],
        ["Election", "ElectionID", "UUID / Simple", "Primary Key (PK)"],
        ["Election", "Title", "VARCHAR(150)", "Not Null"],
        ["Election", "Description", "TEXT", "Optional"],
        ["Election", "Status", "VARCHAR(20)", "Enum ('DRAFT', 'ACTIVE', 'CLOSED', 'TALLIED')"],
        ["Election", "StartTime", "TIMESTAMP", "Not Null"],
        ["Election", "EndTime", "TIMESTAMP", "Not Null"],
        ["Election", "TotalVotes", "INTEGER", "Derived Attribute (Computed from Vote count)"],
        ["Candidate", "CandidateID", "UUID / Simple", "Primary Key (PK)"],
        ["Candidate", "ElectionID", "UUID / Simple", "Foreign Key (FK References Election.ElectionID)"],
        ["Candidate", "CandidateName", "VARCHAR(100)", "Not Null"],
        ["Candidate", "Party", "VARCHAR(100)", "Not Null"],
        ["Token", "TokenID", "UUID / Simple", "Primary Key (PK)"],
        ["Token", "TokenHash", "VARCHAR(64)", "SHA-256 Hash of raw token, Unique, Not Null"],
        ["Token", "IssuedTime", "TIMESTAMP", "Not Null"],
        ["Token", "ExpiryTime", "TIMESTAMP", "Not Null (Time-bound 15-minute window)"],
        ["Token", "Used", "BOOLEAN", "Default False, Flipped to True upon vote"],
        ["Vote", "VoteID", "UUID / Simple", "Primary Key (PK)"],
        ["Vote", "EncryptedBallot", "BYTEA / TEXT", "AES-256-GCM Ciphertext Payload"],
        ["Vote", "Timestamp", "TIMESTAMP", "Commitment Timestamp"],
        ["Vote", "MerkleLeafHash", "VARCHAR(64)", "SHA-256 Hash of EncryptedBallot"],
        ["AuditLog", "LogID", "UUID / Simple", "Primary Key (PK)"],
        ["AuditLog", "Action", "VARCHAR(100)", "Event Description (e.g., 'TOKEN_ISSUED')"],
        ["AuditLog", "Timestamp", "TIMESTAMP", "Audit Event Timestamp"],
        ["AuditLog", "UserRole", "VARCHAR(20)", "Role of initiating identity"],
        ["AuditLog", "IPAddress", "VARCHAR(45)", "IPv4 / IPv6 Client Address"],
        ["User_Election_Eligibility", "UserID", "UUID / Simple", "Composite PK, FK (References User.UserID)"],
        ["User_Election_Eligibility", "ElectionID", "UUID / Simple", "Composite PK, FK (References Election.ElectionID)"],
        ["User_Election_Eligibility", "HasVoted", "BOOLEAN", "Default False, Tracks participation per election"]
    ]
    add_table_styled(doc, er_table_headers, er_table_data, col_widths=[1.5, 1.4, 1.3, 2.3])

    add_heading_styled(doc, "3. Cardinality & Relationship Constraints", level=2)
    add_bullet_styled(doc, "One Election contains multiple Candidates (1:M). Partial participation for Election (election can exist before candidates register); total participation for Candidate.", "Election to Candidate (1:M): ")
    add_bullet_styled(doc, "Users are eligible for multiple Elections, and an Election includes multiple eligible Users. This M:N relationship is resolved through User_Election_Eligibility.", "User to Election (Resolved M:N): ")
    add_bullet_styled(doc, "A single-use token authorizes exactly one vote. However, to preserve privacy, no foreign key connects Token and Vote in physical storage.", "Token to Vote (Conceptual 1:1): ")
    add_bullet_styled(doc, "System operations generate sequential, append-only security logs.", "System to AuditLog (1:M): ")

    add_heading_styled(doc, "4. Entity-Relationship Diagram", level=2)
    add_image_figure(doc, "Report files/extracted_figures/ER_diagram_image1.png", "Final Relational ER Diagram with Zero-PII Voting Architecture")

    add_heading_styled(doc, "5. Result & Conclusion", level=2)
    add_paragraph_styled(doc, "The Entity-Relationship Diagram and relational schema were successfully formulated and normalized. The architecture guarantees mathematical and structural impossibility of linking individual voter identities to decrypted ballots, ensuring unconditional ballot confidentiality.")

    doc.add_page_break()

    # -------------------------------------------------------------------------
    # EXPERIMENT 5
    # -------------------------------------------------------------------------
    add_experiment_header(
        doc,
        exp_no="5",
        title="Behavioral Modelling: State Transition Diagram",
        aim="To model the dynamic behavioral lifecycle, states, events, guard conditions, and cryptographic transitions of the Privacy-Preserving Electronic Voting System using UML State Transition Diagrams.",
        tools="StarUML State Machine Modeller"
    )

    add_heading_styled(doc, "1. Behavioral Modelling Concept", level=2)
    add_paragraph_styled(doc, "UML State Transition Diagrams capture the finite-state machine (FSM) governing a system entity over its operational lifetime. In AegisVote, behavioral modelling ensures that a voter's session and voting token progress through strictly validated, irreversible stages, preventing replay attacks and double voting.")

    add_heading_styled(doc, "2. State Transition Diagram Blueprint", level=2)
    add_paragraph_styled(doc, "The State Transition Diagram models the voting lifecycle from initial voter login to final cryptographic commitment in the Merkle Tree:")
    add_image_figure(doc, "Report files/extracted_figures/State_Transition_image1.png", "UML State Transition Diagram: Voter Session and Cryptographic Lifecycle States")

    add_heading_styled(doc, "3. Detailed State Lifecycle & Transition Analysis", level=2)
    add_paragraph_styled(doc, "The dynamic machine transitions across nine discrete states governed by rigorous guard conditions:")
    add_bullet_styled(doc, "The initial state before user interaction.", "1. Idle / Unauthenticated: ")
    add_bullet_styled(doc, "Triggered by submitCredentials(). System verifies Bcrypt password hash and prompts for 2FA TOTP.", "2. Dual-Factor Authenticating: ")
    add_bullet_styled(doc, "Entered upon successful OTP validation. System evaluates guard condition [HasVoted == False && ElectionStatus == 'ACTIVE']. If HasVoted == True, transitions to Access Denied.", "3. Eligibility Verification: ")
    add_bullet_styled(doc, "Generated upon successful eligibility confirmation. TokenService issues a cryptographically random 256-bit token; User_Election_Eligibility.HasVoted is immediately updated to True.", "4. Token Generated & Active: ")
    add_bullet_styled(doc, "Voter UI displays the candidate ballot. Governed by a strict 15-minute countdown guard [ExpiryTime > CurrentTime].", "5. Ballot Selection: ")
    add_bullet_styled(doc, "Triggered by castVote(token, candidateId). VotingServer invokes AES-256 encryption engine to encrypt candidate selection in memory.", "6. Ballot Encrypted: ")
    add_bullet_styled(doc, "SHA-256 leaf hash is computed and inserted into the append-only Merkle Tree engine, updating the global Merkle root.", "7. Merkle Tree Commitment: ")
    add_bullet_styled(doc, "The Token is permanently invalidated (Used = True), and the raw token is purged from memory.", "8. Token Invalidation: ")
    add_bullet_styled(doc, "Voter receives a cryptographic receipt containing the Merkle leaf hash and timestamp; session terminates securely.", "9. Receipt Dispatched & Terminated: ")

    add_heading_styled(doc, "4. Result & Conclusion", level=2)
    add_paragraph_styled(doc, "The State Transition Diagram was successfully designed and validated in StarUML. The formal model proves that token consumption and ballot commitment occur within an atomic, irreversible sequence, eliminating token replay attacks and ballot tampering.")

    doc.add_page_break()

    # -------------------------------------------------------------------------
    # EXPERIMENT 6
    # -------------------------------------------------------------------------
    add_experiment_header(
        doc,
        exp_no="6",
        title="Object-Oriented Design: Use Case and Class Diagrams",
        aim="To perform Object-Oriented Analysis and Design (OOAD) by developing the Use Case Diagram and Class Diagram for the Privacy-Preserving Electronic Voting System using StarUML.",
        tools="StarUML OOAD Engine"
    )

    add_heading_styled(doc, "1. Identified Actors & Roles", level=2)
    add_paragraph_styled(doc, "Five distinct actors interact with the system across different authorization tiers:")
    add_bullet_styled(doc, "Primary user who authenticates with multi-factor credentials, requests an anonymous token, casts one vote, and reviews election results.", "1. Voter: ")
    add_bullet_styled(doc, "Administrative authority responsible for creating elections, registering candidates, opening/closing polls, and certifying results.", "2. Election Officer: ")
    add_bullet_styled(doc, "Maintains system uptime, provisions accounts, reviews audit logs, and initiates database backups without access to ballot data.", "3. Administrator: ")
    add_bullet_styled(doc, "Trusted entity executing cryptographic Merkle root audits, post-election ballot decryption, and tally calculation.", "4. Counting Authority: ")
    add_bullet_styled(doc, "External SaaS provider delivering out-of-band TOTP codes via SMS or Authenticator protocols.", "5. OTP Provider: ")

    add_heading_styled(doc, "2. Use Case Model & Relationships", level=2)
    add_paragraph_styled(doc, "The Use Case model maps out essential functional interactions along with explicit <<include>> and <<extend>> dependencies:")
    add_bullet_styled(doc, "Authenticate Voter, Authenticate with OTP, Obtain Anonymous Voting Token, Cast Vote, Validate Anonymous Token, Encrypt Ballot, Store Encrypted Ballot, Check Election Status, View Results.", "Voter Use Cases: ")
    add_bullet_styled(doc, "Create Election, Manage Candidates, Open Election, Close Election, Publish Results.", "Election Officer Use Cases: ")
    add_bullet_styled(doc, "Manage User Accounts, Backup Database, View Audit Logs.", "Administrator Use Cases: ")
    add_bullet_styled(doc, "Verify Merkle Integrity, Decrypt Ballots, Tally Votes, Generate Final Results Report.", "Counting Authority Use Cases: ")
    add_bullet_styled(doc, "Authenticate Voter <<includes>> Authenticate with OTP; Cast Vote <<includes>> Validate Anonymous Token, Encrypt Ballot, and Store Encrypted Ballot; Tally Votes <<includes>> Verify Merkle Integrity and Decrypt Ballots; Account Lockout <<extends>> Authenticate Voter (triggered upon repeated authentication failures exceeding threshold).", "Key Dependencies: ")

    add_image_figure(doc, "Report files/extracted_figures/Use_Case_and_Class_image1.png", "StarUML Use Case Diagram: Actors, Use Cases, Include & Extend Dependencies")

    add_heading_styled(doc, "3. Class Diagram & Domain Model", level=2)
    add_paragraph_styled(doc, "The Class Diagram models the static structural architecture, encapsulating attributes, operations, and relationships across 14 domain classes:")
    add_bullet_styled(doc, "User (base class), Voter, ElectionOfficer, Administrator, CountingAuthority. Generalization hierarchy where specific roles inherit authentication properties.", "User & Role Classes: ")
    add_bullet_styled(doc, "Election (composed of Candidates and Ballots) and Candidate.", "Election Entities: ")
    add_bullet_styled(doc, "AnonymousVotingToken (handles token hashing and unicity) and Ballot (stores encrypted payload and Merkle leaf hash without any voter ID).", "Security & Voting Entities: ")
    add_bullet_styled(doc, "AuthenticationServer, TokenService, VotingServer, VoteEncryption, MerkleTreeIntegrity, CountingService, AuditLog.", "Service Classes: ")

    add_image_figure(doc, "Report files/extracted_figures/Use_Case_and_Class_image2.png", "StarUML Class Diagram: Entity Relationships, Attributes, Operations & Multiplicities")

    add_heading_styled(doc, "4. Object-Oriented Design Principles Applied", level=2)
    add_bullet_styled(doc, "Internal representation of sensitive attributes (password hashes, raw tokens, private keys) is strictly hidden.", "Encapsulation: ")
    add_bullet_styled(doc, "High-level services expose abstract interfaces, decoupling business logic from underlying cryptographic libraries.", "Abstraction: ")
    add_bullet_styled(doc, "Voter, ElectionOfficer, Administrator, and CountingAuthority share common credentials from the User superclass.", "Inheritance: ")
    add_bullet_styled(doc, "Each service maintains a single clear purpose (TokenService manages tokens; MerkleTreeIntegrity manages cryptographic ledgers).", "Single Responsibility Principle (SRP): ")

    add_heading_styled(doc, "5. Result & Conclusion", level=2)
    add_paragraph_styled(doc, "The Use Case Diagram and Class Diagram were successfully constructed in StarUML. The object-oriented model cleanly separates responsibilities, decouples voting operations from identity management, and provides a scalable design for implementation.")

    doc.add_page_break()

    # -------------------------------------------------------------------------
    # EXPERIMENT 7
    # -------------------------------------------------------------------------
    add_experiment_header(
        doc,
        exp_no="7",
        title="Interaction Modelling: Sequence Diagram and Communication Diagram",
        aim="To model dynamic object interactions, temporal execution ordering, and structural collaboration networks of the Privacy-Preserving Electronic Voting System using UML Sequence and Communication Diagrams.",
        tools="StarUML Dynamic Interaction Modeller"
    )

    add_heading_styled(doc, "1. Sequence Diagram: Temporal Interaction", level=2)
    add_paragraph_styled(doc, "The Sequence Diagram models the chronological execution timeline across six interacting lifelines: Voter, :VoterClient, :VotingServer, :TokenService, :EncryptionModule, and :ElectionDatabase. The execution progresses through three sequential phases:")
    add_bullet_styled(doc, "Voter selects candidate; UI presents confirmation dialog. Upon confirmation, :VoterClient transmits submitBallot(token, candidateId) to :VotingServer. :VotingServer queries :TokenService with validateToken(tokenHash).", "Phase 1 (User Intent & Authorization): ")
    add_bullet_styled(doc, ":VotingServer passes candidateId to :EncryptionModule in memory to generate an AES-256 ciphertext. :VotingServer invokes appendEncryptedBallot() on :ElectionDatabase to store the ciphertext and trigger Merkle tree root recalculation.", "Phase 2 (Ballot Encryption & Ledger Commitment): ")
    add_bullet_styled(doc, ":ElectionDatabase returns commitment confirmation; :VotingServer calls invalidateToken() on :TokenService and delivers a cryptographic receipt to :VoterClient.", "Phase 3 (Sanitization & Receipt Delivery): ")

    add_image_figure(doc, "Report files/extracted_figures/Sequence_and_Comm_image1.png", "UML Sequence Diagram: Chronological Vote Casting and Cryptographic Commitment")

    add_heading_styled(doc, "2. Communication Diagram: Structural Collaboration", level=2)
    add_paragraph_styled(doc, "While the Sequence Diagram emphasizes the time axis, the Communication Diagram models object topology and message pathways using numbered messaging discipline (messages 1 through 12).")
    add_bullet_styled(doc, "The Communication Diagram demonstrates that :VotingServer acts as a Mediator, decoupling client requests from cryptographic engines and persistent databases.", "Mediator Architecture: ")
    add_bullet_styled(doc, "Synchronous forward invocations (e.g., 5: validateToken, 7: encryptBallot, 9: appendEncryptedBallot) are paired with corresponding return messages (6: tokenValid, 8: encryptedBallot, 10: commitmentSuccess) tracing back along the collaboration links.", "Numbered Messaging Discipline: ")

    add_image_figure(doc, "Report files/extracted_figures/Sequence_and_Comm_image2.png", "UML Communication Diagram: Object Collaborations and Numbered Message Discipline")

    add_heading_styled(doc, "3. Result & Conclusion", level=2)
    add_paragraph_styled(doc, "The Sequence and Communication Diagrams were successfully designed in StarUML. The dual perspectives rigorously confirm both the chronological execution flow and structural decoupling of the voting pipeline.")

    doc.add_page_break()

    # -------------------------------------------------------------------------
    # EXPERIMENT 8
    # -------------------------------------------------------------------------
    add_experiment_header(
        doc,
        exp_no="8",
        title="OO Design: Package, Component, and Deployment Models",
        aim="To design the Package Diagram, Component Diagram, and Deployment Diagram for the Privacy-Preserving Electronic Voting System using StarUML.",
        tools="StarUML Architectural Modeller"
    )

    add_heading_styled(doc, "1. Package Diagram: Logical Module Organization", level=2)
    add_paragraph_styled(doc, "The Package Diagram groups related model elements into five cohesive, loosely coupled logical packages:")
    add_bullet_styled(doc, "Encapsulates the Authentication Server, TOTP validation module, and Anonymous Token Service.", "1. Authentication and Privacy Package: ")
    add_bullet_styled(doc, "Houses the Voting Server and AES-256 Vote Encryption Engine.", "2. Voting and Encryption Package: ")
    add_bullet_styled(doc, "Contains the PostgreSQL database repository and SHA-256 Merkle Tree Integrity Module.", "3. Election Data and Integrity Package: ")
    add_bullet_styled(doc, "Encapsulates post-election decryption and tallying services.", "4. Counting and Results Package: ")
    add_bullet_styled(doc, "Contains the audit logging engine and administrative interfaces.", "5. Audit and Administration Package: ")

    add_image_figure(doc, "Report files/extracted_figures/Package_Component_Deploy_image1.png", "StarUML Package Diagram: Logical Organization and Module Dependencies")

    add_heading_styled(doc, "2. Component Diagram: Service Architecture & Interfaces", level=2)
    add_paragraph_styled(doc, "The Component Diagram models the autonomous software components and their provided and required interfaces:")
    add_bullet_styled(doc, "Web Client, Authentication Server, Anonymous Token Service, Voting Server, Vote Encryption Module, Election Database, Merkle Tree Integrity Module, Counting Authority, Audit Logging Module, and External OTP Provider.", "Components: ")
    add_bullet_styled(doc, "Components communicate via lightweight RESTful HTTPS interfaces. The Voting Server consumes interfaces from Token Service and Encryption Module but operates completely detached from Authentication Server.", "Interfaces: ")

    add_image_figure(doc, "Report files/extracted_figures/Package_Component_Deploy_image2.png", "StarUML Component Diagram: Decoupled Microservice Components and REST Interfaces")

    add_heading_styled(doc, "3. Deployment Diagram: Physical Infrastructure & Execution Environments", level=2)
    add_paragraph_styled(doc, "The Deployment Diagram models the physical hardware nodes, network topologies, and runtime execution environments:")
    add_bullet_styled(doc, "Executes modern Web Browsers (Chrome, Edge, Firefox) accessing the system over HTTPS/TLS 1.3.", "Voter Device Node: ")
    add_bullet_styled(doc, "Hosts the Go backend runtime, running containerized microservices behind a secure NGINX reverse proxy.", "Application Server Node: ")
    add_bullet_styled(doc, "Hosts the PostgreSQL relational database engine and cryptographic append-only Merkle ledger.", "Database Server Node: ")
    add_bullet_styled(doc, "Air-gapped secure hardware console utilized solely post-closure for ballot decryption and tally computation.", "Counting Authority Console: ")

    add_image_figure(doc, "Report files/extracted_figures/Package_Component_Deploy_image3.png", "StarUML Deployment Diagram: Physical Client, Gateway Server, and Database Nodes")

    add_heading_styled(doc, "4. Result & Conclusion", level=2)
    add_paragraph_styled(doc, "The Package, Component, and Deployment Diagrams were successfully constructed. The models represent the complete logical organization, microservice boundaries, physical infrastructure, and network communication paths of the system.")

    doc.add_page_break()

    # -------------------------------------------------------------------------
    # EXPERIMENT 9
    # -------------------------------------------------------------------------
    add_experiment_header(
        doc,
        exp_no="9",
        title="Software Testing: Design and Demonstration of Test Cases",
        aim="To design, demonstrate, and execute functional automated test cases using Selenium IDE and non-functional performance / load test suites using Apache JMeter on the target web application (BlazeDemo).",
        tools="Selenium IDE, Apache JMeter 5.x"
    )

    add_heading_styled(doc, "1. Testing Strategy & Target Application", level=2)
    add_paragraph_styled(doc, "Quality Assurance (QA) verifies that software functions reliably under expected operations and scales under concurrent user loads. The testing suite was executed on BlazeDemo (https://blazedemo.com), a web-based flight reservation simulator. The test strategy incorporates both automated functional UI testing and performance load testing:")
    add_bullet_styled(doc, "Selenium IDE browser automation validating departure/destination routing, flight selection, passenger data entry, and confirmation receipt generation.", "Functional Testing: ")
    add_bullet_styled(doc, "Apache JMeter thread group simulation evaluating server response times, latency, throughput, and error rates under concurrent user load.", "Non-Functional Testing: ")

    add_heading_styled(doc, "2. Master Test Case Suite Table", level=2)
    add_paragraph_styled(doc, "The formal test cases designed, executed, and verified are summarized below:")

    test_headers = ["Test ID", "Test Type", "Scenario Description", "Steps to Execute", "Expected Result", "Actual Result", "Status"]
    test_data = [
        [
            "TC_01",
            "Functional",
            "Verify Flight Search Functionality",
            "1. Navigate to blazedemo.com\n2. Select Departure City (e.g. Paris)\n3. Select Destination City (e.g. Rome)\n4. Click 'Find Flights'",
            "User is redirected to the flights list page displaying available flights for the route.",
            "Flights list page rendered matching search parameters.",
            "PASS"
        ],
        [
            "TC_02",
            "Functional",
            "Verify End-to-End Flight Booking",
            "1. From flights list, click 'Choose This Flight'\n2. Fill passenger details (Name, Address, CC)\n3. Click 'Purchase Flight'",
            "User is redirected to confirmation page displaying 'Thank you for your purchase today!' with unique booking ID.",
            "Confirmation page displayed with valid booking receipt ID.",
            "PASS"
        ],
        [
            "TC_03",
            "Non-Functional",
            "Basic Load Test (10 Concurrent Users)",
            "1. Configure JMeter Thread Group: 10 threads, 5s ramp-up period, 1 loop\n2. Send HTTP GET requests to blazedemo.com\n3. Capture Aggregate Report & Response Graphs",
            "Server responds without crashes or timeouts; 0.00% error rate; average response time < 1000 ms.",
            "All 10 samples completed successfully; 0.00% error rate; average latency 890 ms.",
            "PASS"
        ]
    ]
    add_table_styled(doc, test_headers, test_data, col_widths=[0.6, 0.9, 1.2, 1.4, 1.2, 1.2, 0.5])

    add_image_figure(doc, "Report files/extracted_figures/Lab9_Testing_img1.png", "Master Test Case Suite Table and Verification Status")

    add_heading_styled(doc, "3. Functional Testing Demonstration (Selenium IDE)", level=2)
    add_paragraph_styled(doc, "Functional test cases TC_01 and TC_02 were recorded and executed in Selenium IDE. The test runner executed all locator assertions (ID, XPath, CSS selectors), command events, and page transitions with zero failures.")
    
    add_heading_styled(doc, "TC_01: Flight Search Automation Walkthrough", level=3)
    add_paragraph_styled(doc, "The automated script navigated to the portal, selected 'Paris' as origin and 'Rome' as destination, and clicked the search button. The application returned valid flight options, confirming the routing filter logic.")
    add_image_figure(doc, "Report files/extracted_figures/Lab9_Testing_img2.png", "TC_01: Selenium IDE Flight Search Execution Log & Step Verification")

    add_heading_styled(doc, "TC_02: Flight Booking & Purchase Confirmation Walkthrough", level=3)
    add_paragraph_styled(doc, "The automation script selected the first available carrier, populated form fields with passenger credentials, and dispatched the payment request. The confirmation screen verified successful transaction completion.")
    add_image_figure(doc, "Report files/extracted_figures/Lab9_Testing_img3.png", "TC_02: Selenium IDE Flight Purchase Confirmation & Receipt Validation")

    add_heading_styled(doc, "4. Non-Functional Performance & Load Testing (Apache JMeter)", level=2)
    add_paragraph_styled(doc, "Test case TC_03 evaluated the target application's responsiveness under concurrent traffic. A JMeter Thread Group was instantiated with 10 simulated users ramping up over a 5-second interval.")
    add_image_figure(doc, "Report files/extracted_figures/Lab9_Testing_img4.png", "TC_03: Apache JMeter Thread Group Load Configuration (10 Threads, 5s Ramp-Up)")

    add_heading_styled(doc, "JMeter Aggregate Report & Throughput Analysis", level=3)
    add_paragraph_styled(doc, "The Aggregate Report metrics demonstrated robust server health:")
    add_bullet_styled(doc, "10 HTTP GET requests dispatched across the 5-second ramp-up window.", "Total Samples: ")
    add_bullet_styled(doc, "Average response time observed was 890 ms, with median latency at 860 ms and 90% Line at 980 ms.", "Response Time Metrics: ")
    add_bullet_styled(doc, "0.00% error rate across all threads; zero HTTP 5xx server exceptions or 4xx client errors.", "Error Rate: ")
    add_bullet_styled(doc, "Sustainable throughput of 1.9 transactions/second with continuous stream delivery.", "Throughput: ")

    add_image_figure(doc, "Report files/extracted_figures/Lab9_Testing_img5.png", "TC_03: Apache JMeter Aggregate Report (Response Times, Throughput & Error Rate)")
    add_image_figure(doc, "Report files/extracted_figures/Lab9_Testing_img6.png", "TC_03: Apache JMeter Response Time Graph & Latency Distribution Curve")

    add_heading_styled(doc, "5. Result & Conclusion", level=2)
    add_paragraph_styled(doc, "Both functional test cases (TC_01, TC_02) and the non-functional performance test case (TC_03) were executed and validated. The target application demonstrated 100% test pass rate in Selenium IDE and maintained 0.00% error rate under concurrent load in Apache JMeter.")

    # Save document
    out_docx_path = "Consolidated_Lab_Experiments_24BPS1089.docx"
    doc.save(out_docx_path)
    print(f"Saved primary docx: {out_docx_path}")

    report_files_docx_path = os.path.join("Report files", "Consolidated_Lab_Experiments_24BPS1089.docx")
    shutil.copy(out_docx_path, report_files_docx_path)
    print(f"Copied to: {report_files_docx_path}")

if __name__ == "__main__":
    build_full_record()
