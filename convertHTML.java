import java.io.File;
import java.io.FileInputStream;
import java.io.InputStreamReader;
//w ww  .  j a va  2  s.  co m
import javax.xml.parsers.DocumentBuilder;
import javax.xml.parsers.DocumentBuilderFactory;
import javax.xml.transform.OutputKeys;
import javax.xml.transform.Result;
import javax.xml.transform.Source;
import javax.xml.transform.Transformer;
import javax.xml.transform.TransformerFactory;
import javax.xml.transform.dom.DOMSource;
import javax.xml.transform.stream.StreamResult;

import org.w3c.dom.Document;
import org.xml.sax.InputSource;

public class convertHTML {
  public static void main(String[] argv) throws Exception {
    DocumentBuilderFactory factory = DocumentBuilderFactory.newInstance();
    DocumentBuilder builder = factory.newDocumentBuilder();
    Document document = builder.parse(new InputSource(new InputStreamReader(new FileInputStream(
        "result.xml"))));

    Transformer xformer = TransformerFactory.newInstance().newTransformer();
    xformer.setOutputProperty(OutputKeys.METHOD, "xml");
    xformer.setOutputProperty(OutputKeys.INDENT, "yes");
    //xformer.setOutputProperty("http://xml.apache.org/xslt;indent-amount", "4");
    xformer.setOutputProperty(OutputKeys.OMIT_XML_DECLARATION, "yes");
    Source source = new DOMSource(document);
    Result result = new StreamResult(new File("result2.xslt"));
    xformer.transform(source, result);
  }
}